package main

import (
	"flag"
	"fmt"
	"os"
	importpath "path"
	"path/filepath"
	"strings"

	"server_gin/config"

	"gorm.io/driver/mysql"
	"gorm.io/gen"
	"gorm.io/gorm"
)

const defaultOutPath = "./dao/query"

func main() {
	if err := run(); err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}
}

func run() error {
	var dsn string
	var channel string
	var outPath string
	var modelPkg string
	var tables string
	var withQuery bool

	flag.StringVar(&dsn, "dsn", strings.TrimSpace(os.Getenv("GEN_DSN")), "MySQL DSN. If empty, read from config database.<channel>.link")
	flag.StringVar(&channel, "channel", "default", "database channel name from config")
	flag.StringVar(&outPath, "out", defaultOutPath, "generated query output path")
	flag.StringVar(&modelPkg, "modelPkg", "model", "generated model package name under query output path")
	flag.StringVar(&tables, "tables", "", "comma-separated table names. Empty means all tables")
	flag.BoolVar(&withQuery, "query", true, "generate type-safe query code")
	flag.Parse()

	if dsn == "" {
		cfg, err := config.Load()
		if err != nil {
			return err
		}
		dbCfg, ok := cfg.Database[channel]
		if !ok {
			return fmt.Errorf("database channel %q not found", channel)
		}
		dsn = strings.TrimPrefix(dbCfg.Link, "mysql:")
	}

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		return fmt.Errorf("open db: %w", err)
	}

	g := gen.NewGenerator(gen.Config{
		OutPath:           outPath,
		ModelPkgPath:      modelPkg,
		Mode:              gen.WithoutContext | gen.WithDefaultQuery | gen.WithQueryInterface,
		FieldWithIndexTag: true,
		FieldWithTypeTag:  true,
	})
	g.UseDB(db)

	models, err := generateModels(g, tables)
	if err != nil {
		return err
	}
	if withQuery {
		g.ApplyBasic(models...)
	}
	g.Execute()

	if err := moveGeneratedModel(outPath, modelPkg); err != nil {
		return err
	}

	fmt.Printf("generated %d model(s) into %s\n", len(models), outPath)
	return nil
}

func generateModels(g *gen.Generator, tables string) ([]any, error) {
	items := splitCSV(tables)
	if len(items) == 0 {
		return g.GenerateAllTable(), nil
	}

	models := make([]any, 0, len(items))
	for _, table := range items {
		models = append(models, g.GenerateModel(table))
	}
	return models, nil
}

func moveGeneratedModel(outPath, modelPkg string) error {
	modelPkgPath := filepath.Clean(filepath.FromSlash(modelPkg))
	generatedModelPath := filepath.Join(outPath, modelPkgPath)
	if _, err := os.Stat(generatedModelPath); err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return err
	}

	targetModelPath := filepath.Clean(filepath.Join(outPath, "..", filepath.Base(modelPkgPath)))
	if err := os.RemoveAll(targetModelPath); err != nil {
		return err
	}
	if err := os.Rename(generatedModelPath, targetModelPath); err != nil {
		return err
	}
	return rewriteQueryModelImports(outPath, modelPkgPath, targetModelPath)
}

func rewriteQueryModelImports(outPath, modelPkgPath, targetModelPath string) error {
	module, err := readModulePath()
	if err != nil {
		return err
	}
	fromImport := importpath.Join(module, filepath.ToSlash(filepath.Clean(outPath)), filepath.ToSlash(modelPkgPath))
	toImport := importpath.Join(module, filepath.ToSlash(targetModelPath))

	return filepath.WalkDir(outPath, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() || !strings.HasSuffix(path, ".go") {
			return nil
		}
		b, err := os.ReadFile(path)
		if err != nil {
			return err
		}
		s := strings.ReplaceAll(string(b), fromImport, toImport)
		return os.WriteFile(path, []byte(s), 0644)
	})
}

func readModulePath() (string, error) {
	b, err := os.ReadFile("go.mod")
	if err != nil {
		return "", fmt.Errorf("read go.mod: %w", err)
	}
	for _, line := range strings.Split(string(b), "\n") {
		line = strings.TrimSpace(line)
		if strings.HasPrefix(line, "module ") {
			module := strings.TrimSpace(strings.TrimPrefix(line, "module "))
			if module != "" {
				return module, nil
			}
		}
	}
	return "", fmt.Errorf("module path not found in go.mod")
}

func splitCSV(s string) []string {
	parts := strings.Split(s, ",")
	out := make([]string, 0, len(parts))
	for _, part := range parts {
		v := strings.TrimSpace(part)
		if v != "" {
			out = append(out, v)
		}
	}
	return out
}

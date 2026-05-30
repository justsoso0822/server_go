package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"sort"
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

	if err := moveGeneratedModel(outPath); err != nil {
		return err
	}

	fmt.Printf("generated %d model(s) into %s\n", len(models), outPath)
	return nil
}

func generateModels(g *gen.Generator, tables string) ([]any, error) {
	items := splitCSV(tables)
	if len(items) == 0 {
		items = knownTables()
	}

	models := make([]any, 0, len(items))
	for _, table := range items {
		if name, ok := modelNameByTable[table]; ok {
			models = append(models, g.GenerateModelAs(table, name))
			continue
		}
		models = append(models, g.GenerateModel(table))
	}
	return models, nil
}

var modelNameByTable = map[string]string{
	"log_login":        "LogLogin",
	"log_msg":          "LogMsg",
	"log_trace":        "LogTrace",
	"mem_config":       "MemConfig",
	"prf_flower":       "PrfFlower",
	"prf_flower_level": "PrfFlowerLevel",
	"prf_item":         "PrfItem",
	"prf_res":          "PrfRes",
	"prf_task":         "PrfTask",
	"sys_die":          "SysDie",
	"sys_gm":           "SysGm",
	"user":             "User",
	"user_bag":         "UserBag",
	"user_bag_tp":      "UserBagTp",
	"user_data":        "UserData",
	"user_item":        "UserItem",
	"user_loginkey":    "UserLoginkey",
	"user_online":      "UserOnline",
	"user_res":         "UserRes",
	"user_task":        "UserTask",
}

func knownTables() []string {
	out := make([]string, 0, len(modelNameByTable))
	for table := range modelNameByTable {
		out = append(out, table)
	}
	sort.Strings(out)
	return out
}

func moveGeneratedModel(outPath string) error {
	generatedModelPath := filepath.Join(outPath, "model")
	if _, err := os.Stat(generatedModelPath); err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return err
	}

	targetModelPath := filepath.Clean(filepath.Join(outPath, "..", "model"))
	if err := os.RemoveAll(targetModelPath); err != nil {
		return err
	}
	if err := os.Rename(generatedModelPath, targetModelPath); err != nil {
		return err
	}
	return rewriteQueryModelImports(outPath)
}

func rewriteQueryModelImports(outPath string) error {
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
		s := strings.ReplaceAll(string(b), "server_gin/dao/query/model", "server_gin/dao/model")
		return os.WriteFile(path, []byte(s), 0644)
	})
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

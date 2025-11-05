package app

import (
	"context"
	"fmt"
	"github.com/CyberGeo335/prak_five/internal/utils"
	"github.com/joho/godotenv"
	"log"
	"time"
)

func Run() {

	_ = godotenv.Load()

	built, err := utils.BuildPostgresURL()
	if err != nil {
		log.Fatalf("Error building Postgres URL: %s", err)
	}
	db, err := utils.OpenDB(built)
	if err != nil {
		log.Fatalf("Error opening DB: %s", err)
	}
	defer db.Close()

	repo := utils.NewRepo(db)

	// 1) Вставим пару задач
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	titles := []string{"Сделать ПЗ №5", "Купить кофе", "Проверить отчёты"}
	for _, title := range titles {
		id, err := repo.CreateTask(ctx, title)
		if err != nil {
			log.Fatalf("CreateTask error: %v", err)
		}
		log.Printf("Inserted task id=%d (%s)", id, title)
	}

	// 2) Прочитаем список задач
	ctxList, cancelList := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancelList()

	tasks, err := repo.ListTasks(ctxList)
	if err != nil {
		log.Fatalf("ListTasks error: %v", err)
	}

	// 3) Напечатаем
	fmt.Println("=== Tasks ===")
	for _, t := range tasks {
		fmt.Printf("#%d | %-24s | done=%-5v | %s\n",
			t.ID, t.Title, t.Done, t.CreatedAt.Format(time.RFC3339))
	}
}

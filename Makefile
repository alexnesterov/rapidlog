MIN_COVERAGE = 37

.PHONY: all build test test-fast test-cov test-html test-check test-tdd

all: build test
	@echo "🎉 Все проверки пройдены!"

build:
	@echo "🏗️  Проверка сборки всех пакетов (go build ./...)..."
	go build ./...
	@echo "✅ Сборка успешно проверена!"

test:
	@echo "🧪 Запуск всех тестов проекта..."
	@go test ./... && echo "🟩 [GREEN] Отлично! Все тесты прошли." || (echo "🟥 [RED] Тесты упали! Исправляй код." && exit 1)

test-fast:
	@echo "⚡ Быстрый перезапуск тестов без кэша..."
	@go test -count=1 ./... && echo "🟩 [GREEN] Отлично! Все тесты прошли." || (echo "🟥 [RED] Тесты упали! Исправляй код." && exit 1)

test-cov:
	@echo "🧪 Запуск тестов со сквозным покрытием пакетов..."
	go test -coverprofile=coverage.out -coverpkg=./... ./...

test-html: test-cov
	@echo "🌐 Открытие отчета в браузере..."
	go tool cover -html=coverage.out

test-check: test-cov
	@echo "🛡️ Проверка минимального порога покрытия ($(MIN_COVERAGE)%)..."
	@TOTAL_COV=$$(go tool cover -func=coverage.out | grep "total:" | awk '{print $$3}' | sed 's/%//'); \
	echo "📊 Текущее общее покрытие: $$TOTAL_COV%"; \
	BC_CHECK=$$(echo "$$TOTAL_COV >= $(MIN_COVERAGE)" | bc -l 2>/dev/null || awk "BEGIN {if ($$TOTAL_COV >= $(MIN_COVERAGE)) print 1; else print 0}"); \
	if [ "$$BC_CHECK" -eq 1 ]; then \
		echo "✅ Проверка пройдена! Покрытие в норме."; \
	else \
		echo "❌ Ошибка: Общее покрытие ($$TOTAL_COV%) ниже допустимого порога ($(MIN_COVERAGE)%)!"; \
		exit 1; \
	fi

dir ?= ./...

test-tdd:
	@echo "🔄 [TDD] Запуск тестов в $(dir)..."
	@go test -v -count=1 $(dir) && echo "🟩 [GREEN] Отлично! Все тесты прошли." || (echo "🟥 [RED] Тесты упали! Исправляй код." && exit 1)

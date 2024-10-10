schema-docs:
	@echo "Generating schema documentation..."
	go run utils/gen_schema.go
	@echo "Schema documentation generated."

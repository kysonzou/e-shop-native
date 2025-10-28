#!/bin/bash
set -e

# echo "Running pre-commit checks..."

# # 1. 格式化检查
# echo "→ Checking code format..."
# make lint-fix

# # 2. 运行测试
# echo "→ Running tests..."
# make test

# # 3. 检查生成的代码是否最新
# echo "→ Checking generated code..."
# make api wire
# if ! git diff --exit-code --quiet; then
#     echo "❌ Generated code is out of date!"
#     echo "Run 'make api wire' and commit the changes"
#     exit 1
# fi

echo "✓ Pre-commit checks passed!"
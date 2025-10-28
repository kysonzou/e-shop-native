#!/bin/bash

# 获取 commit message 文件路径（Git 会传入）
commit_msg_file="$1"

# 如果没有参数，从 stdin 读取（用于测试）
if [ -z "$commit_msg_file" ]; then
    commit_msg=$(cat)
else
    commit_msg=$(cat "$commit_msg_file")
fi

# 允许的类型
types="feat|fix|docs|style|refactor|test|chore|perf|ci|build|revert"

# 检查格式
if ! echo "$commit_msg" | grep -qE "^($types)(\(.+\))?: .{1,}"; then
    echo ""
    echo "❌ Invalid commit message format!"
    echo ""
    echo "📝 Format: <type>(<scope>): <subject>"
    echo ""
    echo "📌 Types:"
    echo "  feat:     ✨ New feature"
    echo "  fix:      🐛 Bug fix"
    echo "  docs:     📚 Documentation"
    echo "  style:    💄 Code style"
    echo "  refactor: ♻️  Refactoring"
    echo "  test:     ✅ Tests"
    echo "  chore:    🔧 Maintenance"
    echo "  perf:     ⚡ Performance"
    echo "  ci:       👷 CI/CD"
    echo "  build:    📦 Build"
    echo "  revert:   ⏪ Revert"
    echo ""
    echo "✏️  Examples:"
    echo "  feat(user): add login feature"
    echo "  fix(api): correct response format"
    echo "  docs: update README"
    echo ""
    echo "Your message was: '$commit_msg'"
    echo ""
    exit 1
fi

echo "✓ Commit message format is valid"
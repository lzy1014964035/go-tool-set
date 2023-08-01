#!/bin/bash

# git add .
# git commit -m "自动提交"

# 删除标签
git tag -d v0.0.2
# 重新打标签
git tag -a v0.0.2 -m "开发用标签，每次推送时会自动打这个标签"
# 删除仓库标签
git push --delete origin v0.0.2
# 推送标签至仓库
git push origin v0.0.2
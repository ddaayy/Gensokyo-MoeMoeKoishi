# Changelog — Release008

> 自 Release007 (`a44882c`) 以来的所有变更。

---

## 🚀 新增功能

### CQ:file 文件上传支持

QQ 机器人 API v2 现已开放 `file_type=4`（文件）的富媒体上传，可通过 `POST /v2/users/{openid}/files` 和 `POST /v2/groups/{group_openid}/files` 上传任意文件。

Gensokyo 新增 `[CQ:file]` CQ 码的完整支持：

- `[CQ:file,file=file:///path/to/file]` — 本地文件路径（自动 base64 编码后走 CDN 上传）
- `[CQ:file,file=http://example.com/file.zip]` — HTTP 远程文件链接
- `[CQ:file,file=https://example.com/file.zip]` — HTTPS 远程文件链接
- `[CQ:file,file=base64://<data>]` — base64 编码数据（走 CDN 上传）

支持可选参数 `file_name` 指定文件名，预留给未来 API 开放使用：
- `[CQ:file,file=file:///path/to/file,file_name=myfile.txt]`
- 数组段格式：`{"type":"file","data":{"file":"file:///path","file_name":"myfile.txt"}}`
- 不填时自动从路径/URL 末尾提取（`filepath.Base()`）

支持场景：
- 群聊发送文件（`send_group_msg`）
- 私聊发送文件（`send_private_msg`）
- 私聊互动召回消息（`send_private_msg_wakeup`）

> 💡 **文件名：** 使用 `file_data`（base64）上传时，`file_name` 参数会直接写入 QQ Bot API 的 `file_name` 字段使文件正确命名。URL 方式发送时 QQ 默认从 URL 末尾提取文件名，也可用 `file_name` 覆盖。

### send_private_msg_wakeup API

新增 `send_private_msg_wakeup` API，用于向 QQ 用户发送 C2C 互动召回（唤醒）消息。OneBot 应用端可通过此接口主动唤醒用户会话，不受被动回复上下文限制。

---

## 🔧 改进

### send_private_msg_wakeup 处理优化

- 添加 `active` / `active_type` / `active_sub_type` key 遍历跳过逻辑，避免 `[CQ:active]` 内容被错误地作为媒体消息发送
- 纯 `[CQ:active]` 无实际内容时发送空白唤醒请求，确保用户收到互动通知
- 调用方获得真实的成功/失败返回（同步模式），不再因异步处理导致 WebSocket 超时

---

## 🐛 Bug 修复

### 语音 URL 未正确重命名引起的潜在 panic

**文件：** `handlers/send_group_msg.go`

`url_record` 分支中 `generateGroupMessage` 使用了外层作用域的 `imageURLs` 变量而非 `recordURLs`。编译通过但在实际语音发送路径下会因索引越界 panic。已修正为 `recordURLs[0]`。

### send_private_msg_wakeup 遍历 active key 时误发媒体

`foundItems` 遍历时未跳过 `"active"` key，导致 `[CQ:active]` 标记被当作媒体类型发送，产生 `"Expected RichMediaMessage type for key active"` 错误。已添加 key 过滤。

### NoneBot 插件 msg_text 截断

**文件：** `active_msg/__init__.py`

`handle_wakeup` 中 `on_command("唤醒")` 已将命令前缀剥离，但代码中仍使用 `text[len("唤醒"):].strip()` 再次裁剪，导致 `"唤醒 @target 123"` 的实际消息变为 `"3"`。已修复为直接使用 `text`。

### 富媒体消息 FileType 注释更新

`botgo/dto/message_create.go` 中 `RichMediaMessage.FileType` 注释新增 `4 文件` 类型。

### 文件消息段未处理导致静默丢弃

**文件：** `handlers/message_parser.go`

NoneBot 以 koishi 数组段格式 `{"type":"file","data":{"file":"file:///..."}}` 发送文件消息时，`parseMessageContent` 的 `switch segmentType` 中没有 `case "file":`，日志打印 `Unhandled segment type: file`，文件被静默丢弃。已在 koishi 和 TRSS 两种消息格式中均添加 `case "file":` 处理。

### 本地文件路径 URL 编码未解码

**文件：** `handlers/message_parser.go`

`file:///` 路径中的中文等字符经 URL 编码（如 `%E7%A5%9E` → `神`），去掉 `file:///` 前缀后路径仍为编码状态，`os.ReadFile` 找不到文件。已在两处 `case "file":` 中添加 `neturl.PathUnescape()` 解码。

### foundItems 遍历缺少文件类型 key

**文件：** `handlers/send_group_msg.go`、`handlers/send_private_msg.go`

`local_file` / `base64_file` 经 `generateGroupMessage` 上传 CDN 后返回 `MessageToCreate`，但遍历 `foundItems` 时 `keyMap` 中没有这些 key，导致上传成功的文件不会被发送。已添加 `local_file`、`url_file`、`url_files`、`base64_file` 到 `keyMap`。

---

## 📦 文件变更清单

| 文件 | 变更 |
|------|------|
| `botgo/dto/message_create.go` | FileType 注释新增 `4 文件` 类型 |
| `handlers/message_parser.go` | 新增 CQ:file 正则解析 + foundItems key 映射 + 数组段 `case "file"` + URL 解码 |
| `handlers/send_group_msg.go` | `generateGroupMessage`/`generatePrivateMessage` 文件处理分支 + keyMap 补充 + 文件名传递 |
| `handlers/send_private_msg.go` | keyMap 补充文件类型 + RichMediaMessage 上传后文件名透传 |
| `handlers/send_private_msg_wakeup.go` | 同步模式改造 + active key 跳过 + 空内容兜底 |

---

## 🚀 新增功能（续）

### 统一图床包 `imagehosting`（oss_type 4~10）

新增 `imagehosting/` 统一图床包，提供 7 种后端：

| 后端 | oss_type | 说明 |
|------|----------|------|
| COS 自签 | 4 | 腾讯云 COS（HMAC-SHA1 自签，无需 SDK） |
| Bilibili | 5 | B站开放平台图片上传（需 Cookie） |
| QQ频道 | 6 | 通过发消息获取 qpic.cn 链接 |
| ChatGLM | 7 | 智谱免费图床，开箱即用 |
| Ukaka | 8 | 免费图床，开箱即用 |
| 星野 | 9 | 免费图床，开箱即用 |
| Nature | 10 | 腾讯 COS 直传（密钥内置），开箱即用 |

- `config/config.go` 新增 `OssTypeCOS` ~ `OssTypeNature` 常量（4~10）
- `imagehosting/hosting.go` 提供 `UploadProvider(name, data, filename)` 统一入口
- `images/upload_api.go` 中 `UploadBase64ImageToServer` 按 `oss_type` 分发到对应后端
- `template/config_template.go` 新增 `image_hosting` 示例段

### 配置模板 `image_hosting` 段

`config_template.go` 新增完整的 `image_hosting` 配置段，涵盖 7 种后端的示例字段。

---

## 🔧 改进（续）

### 配置自动补全增强

- **YAML 完整块提取** — `extractMissingConfigLines` 现提取完整的 YAML 块（含注释和子字段），而非单行，避免补全后 yaml 结构错误
- **跳过已存在父块的子 key** — `extractMissingConfigLines` 检测到父块已在配置中时，跳过其子 key 的补全，防止重复插入
- **祖先追溯逻辑** — `buildParentKeyMap` 支持多层祖先追溯，确保嵌套配置补全位置正确

### 文档完善

- `readme.md` 功能列表/CQ 码/API/鸣谢按最新状态更新，移除 "todo,正在施工..."
- `template/config_template.go` 中 `text_intent` 按模板顺序排列，补全所有 intents 列表
- `image_hosting` 注释完善，标注 oss_type 对应关系

---

## 📝 文档

- `readme.md` — 功能列表新增 `[CQ:file]`、`send_private_msg_wakeup`、imagehosting 等；CQ 码/API/Event 表格同步最新实现；配置示例替换为完整可用版本；鸣谢更新
- `docs/cq码/标准CQ码/标准cq码-cq-file.md` — 新增 CQ:file 使用文档
- `release_log/CHANGELOG_v007.md` — 补充 `[CQ:active]` 条目
- `release_log/CHANGELOG_v008.md` — 本文档

---

## 🐛 Bug 修复（续）

### Markdown 图片本地文件路径重写

**文件：** `handlers/message_parser.go`

Markdown 内容中的 `![](本地路径)` 图片在自动上传到 QQ CDN 后，仅替换了 `http://` 和 `https://` 的 URL 模式，未正确处理 `file:///` 协议和纯本地路径。修复后三种路径均能正确识别并替换为 CDN 直链。

### 文件 URL 未正确替换

**文件：** `handlers/message_parser.go`

部分语法中的文件 URL（如 `file:///C:/path/image.png`）在重写 Markdown 图片时未被正确识别和替换，导致图片显示为空白。已修复正则匹配逻辑。

---

## 📦 文件变更清单（续）

| 文件 | 变更 |
|------|------|
| `config/config.go` | 新增 OssTypeCOS ~ OssTypeNature 常量（4~10） |
| `template/config_template.go` | 新增 `image_hosting` 配置段；text_intent 按序排列 |
| `structs/structs.go` | 新增 ImageHostingConfig 及子结构体 |
| `imagehosting/hosting.go` | 新增，统一调度器 + 辅助函数 |
| `imagehosting/cos.go` | 新增，COS 自签上传 |
| `imagehosting/bilibili.go` | 新增，B站图床 |
| `imagehosting/qq_channel.go` | 新增，QQ频道图床 |
| `imagehosting/chatglm.go` | 新增，智谱免费图床 |
| `imagehosting/signed.go` | 新增，Ukaka + 星野签名上传 |
| `imagehosting/nature.go` | 新增，Nature COS 直传 |
| `imagehosting/utils.go` | 新增，辅助函数 |
| `images/upload_api.go` | 新增 oss_type 4~10 分发逻辑 |
| `handlers/message_parser.go` | Markdown 图片本地路径重写修复 + 文件 URL 替换修复 |
| `readme.md` | 功能列表/CQ码/API/Event/鸣谢全面更新 |

---

## 🚀 新增功能（续）

### 统一图床包 imagehosting

新增 `imagehosting` 包，将免费图床（ChatGLM、Ukaka、星野、Nature）和需凭证的图床（COS 自签、Bilibili、QQ 频道）统一为 `oss_type` 枚举（4~10）的后端实现。用户通过 `oss_type` 选择一个后端，不再配置多个 `enabled`。

- 新增 `UploadProvider(name, data, filename)` 按名称选择单个后端
- 新增 `UploadBase64Provider` / `UploadBytes` 兼容旧接口
- 配置模板新增 `image_hosting` 段，存储各后端凭证

### 配置模板新增 image_hosting 段 / 配置自动补全增强

- 配置模板添加 `image_hosting` 段，按 `oss_type` 填写对应凭证
- `text_intent` 列表按模板顺序排列，补全完整的 11 个 intent
- 配置自动补全重构：提取完整 YAML 块而非单行、跳过已存在父块的子 key、祖先追溯逻辑

---

## 📝 文档

- 精简 readme 鸣谢列表，删除无用鸣谢
- 更新 readme 功能列表、CQ 码、API 表格、Event/Intent 说明
- 替换配置示例为实际可用示例，替换为完整配置示例
- 新增 CQ:file 标准 CQ 码文档
- 完善 image_hosting 注释
- 补全完整 text_intent 列表

---

## 🐛 Bug 修复

### 配置自动补全：提取完整 YAML 块

`extractMissingConfigLines` 逐行提取时遇到多行值（如 `post_*` 数组）只取首行，导致 YAML 结构损坏。已改为提取从 key 行到下一同缩进 key 行之间的完整块。

### 配置自动补全：跳过已存在父块的子 key

`appendToConfigFile` 在祖先块（如 `settings:`）已存在时仍插入其子 key，导致重复。已添加 `keyExistsInConfig` 检查，跳过已存在的父块。

### 配置自动补全：祖先追溯逻辑

`appendToConfigFile` 使用 `lastIndex` 定位插入点，在配置项无显式父块缩进时错误插入到文件末尾。已改为从目标 key 行向上追溯找到父块进行插入。

### Markdown 图片重写：本地文件路径支持

`ResolveMarkdownImages` 中正则 `!\[.*?\]\((.*?)\)` 未匹配本地路径（无协议头），导致本地 Markdown 图片无法上传。已修复正则并添加 `file://` 协议头补充逻辑。

### 部分语法中的文件 URL 无法被正确替换

修复 Markdown 图片重写中的路径解析问题，确保 `file:///` 和纯本地路径均能被正确上传并替换为 CDN 直链。

---

## 📦 文件变更清单（补充）

| 文件 | 变更 |
|------|------|
| `imagehosting/hosting.go` | 新增统一图床调度器 |
| `imagehosting/cos.go` | 腾讯云 COS 自签上传 |
| `imagehosting/bilibili.go` | B站图床上传 |
| `imagehosting/qq_channel.go` | QQ频道图床上传 |
| `imagehosting/chatglm.go` | 智谱免费图床 |
| `imagehosting/signed.go` | Ukaka + 星野签名上传 |
| `imagehosting/nature.go` | Nature 内置密钥 COS 直传 |
| `imagehosting/utils.go` | 辅助函数 |
| `imagehosting/README.md` | 图床文档 |
| `structs/structs.go` | 新增 ImageHostingConfig 结构体 |
| `config/config.go` | 新增 OssType 常量 + GetOssTypeName + GetImageHosting* 访问器 |
| `template/config_template.go` | 新增 image_hosting 配置段 |
| `images/upload_api.go` | 按 oss_type 分发到 imagehosting 后端 |
| `handlers/message_parser.go` | Markdown 图片重写支持本地路径 |

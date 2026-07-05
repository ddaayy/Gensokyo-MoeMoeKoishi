# Changelog — Release007 (since Release006)

> 自 Release006 (`7ec67db`) 以来的所有变更。

---

## 🐛 Bug 修复

### 语音/图片本地图床上传失败（内网 server_dir 场景）

**问题：** 当 `server_dir` 设为内网地址（`127.0.0.1` / `192.168.x.x`）时，语音和图片上传到本地图床后，QQ CDN 拉取不到文件，导致"local server uses a private address"错误，媒体消息发送失败。

**根因：** `send_group_msg.go` 中所有语音和图片上传都走了一条统一的路径：base64 → 上传到 Gensokyo 本地 HTTP 服务器 → 拿到 URL → QQ CDN 拉取。该路径要求 `server_dir` 是公网可访问的地址。

**修复：** 语音和图片（非 Markdown）统一改为直接 base64 上传到 QQ CDN，跳过本地图床中转。

```
修复前: base64 → 本地图床 → 公网URL → QQ CDN 拉取 ✅(公网) ❌(内网)
修复后: base64 → 直接提交 QQ CDN API 上传 ✅(任何网络)
```

仅 Markdown 消息中的图片/视频保留本地图床路径（Markdown 内容需要嵌入公开 URL）。

**涉及修改：**
- `send_group_msg.go` — 移除全部 `GetUploadPicV2Base64()` 条件分支和 `UploadBase64ImageToServer`/`UploadBase64RecordToServer` 调用
- 群聊 `base64_image`、`url_image`、`url_images` → `CreateAndUploadMediaMessage`
- 群聊 `base64_record`、`local_record`、`url_record` → `CreateAndUploadMediaMessage`
- C2C 私聊同路径 → `CreateAndUploadMediaMessagePrivate`
- 共约 -350 行冗余代码

---

## ✅ 提交记录

```
6a59179  语音统一直接base64上传QQ CDN
b68957c  图片统一直接base64上传QQ CDN，仅Markdown用图床
```

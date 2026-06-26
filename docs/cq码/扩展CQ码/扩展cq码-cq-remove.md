# [CQ:remove] — 撤回用户最近消息

## 说明

出站 CQ 码，用于撤回指定用户在**当前群**最近一条消息。

> 与 `delete_msg` API 不同：`[CQ:remove]` 无需显式传入 `message_id`，Gensokyo 内部自动查找目标用户在该群的最近消息 ID 并撤回。

## 格式

```
[CQ:remove,user_id=虚拟用户ID]
```

| 参数 | 类型 | 说明 |
|------|------|------|
| `user_id` | int64 | 目标用户的**虚拟 ID**（Gensokyo 转换后的数字 ID） |

## 流程

```
① 后端插件调用 send_group_msg
   send_group_msg(group_id=821404315, message="[CQ:remove,user_id=3607918353]")

② Gensokyo 解析 CQ 码
   - 剥离 [CQ:remove,...]，不发送到 QQ 频道
   - 虚拟 user_id → 真实 OpenID 反向转换
   - 从 msg 数据库查找该用户在该群的最近一条 message_id（6 分钟 TTL）

③ Gensokyo 调用 QQ API 撤回
   api.RetractGroupMessage(groupOpenID, realMsgID)

④ 若 messageText 剥离后为空，跳过发送（不发空消息）
```

## 限制

| 限制 | 说明 |
|------|------|
| 时效 | 只能撤回 **6 分钟内**的消息（与 QQ 官方撤回窗口对齐，msg 数据库 TTL 设为 6 分钟） |
| 范围 | 仅群聊 |
| 权限 | 需机器人为群管理员（QQ API 要求） |

## 后端示例（nonebot2）

```python
from nonebot.adapters.onebot.v11 import Bot, Message

@on_command("撤回").handle()
async def recall_last(bot: Bot, event, args):
    target_uid = extract_user_id(event, args)  # 获取目标用户的虚拟 ID
    await bot.send_group_msg(
        group_id=event.group_id,
        message=Message(f"[CQ:remove,user_id={target_uid}]")
    )
```

## 内部实现

| 模块 | 职责 |
|------|------|
| `idmap.StoreLatestMsgID` | 群消息入站时记录 `(groupOpenID, userOpenID) → realMsgID` |
| `idmap.GetLatestMsgID` | 出站 CQ 码触发时查询最近消息 ID |
| `idmap.cleanExpiredLatestMsg` | 每分钟清理 ≥6 分钟的过期索引 |
| `handlers.ProcessCQRemoveOutbound` | 解析 `[CQ:remove,...]`，剥离 CQ 码，转换用户 ID |

## 适用范围

🏷️ 群聊（出站）

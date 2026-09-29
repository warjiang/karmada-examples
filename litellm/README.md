# ARK → LiteLLM 模型同步

读取火山方舟（ARK）中状态为 `Running` 的内置托管接入点和用户普通接入点，并把它们同步到 LiteLLM。工具同时查询 `InnerDescribeModelEndpoints`（内置 `ep-m-*`）和 `ListEndpoints`（用户 `ep-*`）；同名时以内置托管接入点为准。

内置托管接入点始终使用完整模型 ID 调用，例如 `volcengine/deepseek-v4-1-flash-260910`，`ep-m-*` 仅作为来源信息保存。用户普通接入点使用 Endpoint ID 调用；缺少可调用 ID 时会计入 `uncallable`。运行结果会分别显示 `non-running`、`uncallable`、`shadowed-custom` 和 `duplicate` 数量。

工具只删除自己创建且带有 `ark_sync_managed=true` 标记的模型，不会改动手工配置的模型。默认仅预览，只有 `--apply` 才会写入。

## 使用

要求 Go 1.22+，并确保 LiteLLM 已连接数据库且设置：

```dotenv
STORE_MODEL_IN_DB=True
```

在本目录下运行（独立 Go 模块 `ark-litellm-sync`，已加入仓库的 `go.work`）。从仓库根目录进入并复制配置、填写凭据；已有 `.env` 时无需重复复制：

```bash
cd litellm
cp .env.example .env
```

其中：

- `VOLCENGINE_ACCESS_KEY_ID` / `VOLCENGINE_SECRET_ACCESS_KEY`：仅用于查询 ARK 控制面。
- `ARK_API_KEY`：写入 LiteLLM，供 `volcengine/<endpoint-id>` 推理使用。
- `VOLCENGINE_PROJECT_NAME`：可选；为空时读取当前账号可见的全部项目。
- `LITELLM_MASTER_KEY`：LiteLLM 管理密钥。

预览：

```bash
make run
```

执行同步：

```bash
make apply
```

运行格式化、测试和静态检查：

```bash
make check
```

建议先预览，再通过 cron 或 CI 定期运行 `--apply`。

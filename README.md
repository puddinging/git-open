# git-open

`git-open` 是一个命令行工具，用于在浏览器中打开当前 Git 仓库的远程 URL。

## 功能特点

- 支持多种 Git 托管平台（如 GitHub, GitLab, Bitbucket 等）
- 自动检测并使用系统默认浏览器
- 支持 Windows, macOS, Linux 以及 WSL 环境

## 安装

### 前置条件

- Go 1.16 或更高版本

### 从源码安装

1. 克隆仓库：
   ```
   git clone https://github.com/your-username/git-open.git
   cd git-open
   ```

2. 编译：
   ```
   go build
   ```

3. 将编译后的二进制文件移动到 PATH 中的目录：

   - macOS 和 Linux:
     ```
     sudo mv git-open /usr/local/bin/
     ```

   - Windows:
     将 `git-open.exe` 移动到 `C:\Windows\System32\` 或其他在 PATH 中的目录

   - WSL:
     ```
     sudo mv git-open /usr/local/bin/
     ```

## 使用方法

在任何 Git 仓库目录中，您可以使用以下命令：

1. 打开当前分支的远程仓库页面：
   ```
   git open
   ```

注意：确保您在 Git 仓库目录中运行这些命令，否则可能会出现错误。

如果您的远程仓库 URL 使用 SSH 协议（如 `git@github.com:user/repo.git`），`git-open` 会自动将其转换为 HTTPS URL 以在浏览器中打开。

### 平台特定说明

- macOS: 默认使用 `open` 命令打开默认浏览器
- Windows: 使用 `powershell.exe Start-Process` 命令打开默认浏览器
- Linux: 使用 `xdg-open` 命令打开默认浏览器
- WSL: 使用 `powershell.exe Start-Process` 命令在 Windows 主机上打开默认浏览器

如果您的系统环境变量中设置了 `BROWSER`，`git-open` 将优先使用该设置打开浏览器。

# GTerm

A terminal tool developed with Wails + Vue3 + TypeScript.

This is the first development preview version of GTerm. Due to my full-time work commitments, I can only dedicate limited time to this project. Development will progress gradually, and currently, it only implements basic functionality. There are still many bugs and issues that need to be addressed. If you're interested in contributing to the project, pull requests are welcome!

## Screenshots

1，终端用dom，不用canvas，windows字体发虚；

2，添加复制粘贴快捷键；

3，文件按linux默认排序；

4，标题栏修改，完整显示，不要居中。

基本的文件右键还没有，后续改成原始js，不用vue。

<table>
<tr>
<td><img width="1200" height="800" alt="image" src="https://github.com/user-attachments/assets/83bd5be2-07fe-4e94-ae48-04ae1817037d" />
</td>
<td><img width="1200" height="800" alt="image" src="https://github.com/user-attachments/assets/92d8e853-b8db-4137-856d-9367d210bfe3" />
</td>
</tr>
</table>

# Supported Platforms

- Windows 10/11 AMD64/ARM64
- MacOS 10.13+ AMD64
- MacOS 11.0+ ARM64
- Linux AMD64/ARM64

## Dependencies

- Go 1.24.1
- PNPM (Node 18+)
- Wails CLI v2.10.1

## Development

```bash
wails dev
```

## Building

```bash
wails build
```

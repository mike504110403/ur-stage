### 相關流程

測個ci
<details>
  <summary>
    會員登入流程
  </summary>
  
``` mermaid
graph TD
    A[會員登入] --> B{帳號是否存在}
    B -->|是| C{密碼是否正確}
    B -->|否| F[登入失敗]
    C -->|是| E[登入成功]
    C -->|是| F
```
</details>

<details>
  <summary>
    會員註冊流程
  </summary>
  
``` mermaid
graph TD
    A[會員註冊] --> G{帳號是否重複}
    G -->|是| H[註冊失敗]
    G -->|否| I{註冊資訊是否合法}
    I -->|是| J[註冊成功]
    I -->|否| H
```
</details>


<details>
    <summary>
        會員狀態排程
    </summary>

``` mermaid
graph TD
	A[排程開始] --> B{有流水異動會員}
	B -->|是| C{是否更新會員VIP等級}
	B -->|否| D[等待]
	C -->|是| F[更新]
	C -->|否| D
	F --> D
	D --> A
```
</details>
<details>
  <summary>
    遊戲帳務流程
  </summary>

``` mermaid
graph TD
A[遊戲帳務Req] --> B[(遊戲帳務sql)]
A --> F[(遊戲帳務Redis)]
A --> C[(注單紀錄)]
```
</details>
test
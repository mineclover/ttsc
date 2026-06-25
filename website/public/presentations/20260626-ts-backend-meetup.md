---
marp: true
theme: default
paginate: true
size: 4:3
title: "TTSC: Transformer Survival in the TypeScript 7 Era"
description: "TypeScript Backend Meetup, 2026-06-26"
footer: "TTSC | TypeScript Backend Meetup | 2026-06-26"
style: |
  section {
    font-family: "Inter", "Segoe UI", "Pretendard", sans-serif;
    color: #111827;
    padding: 72px 84px;
  }
  section.lead {
    display: flex;
    flex-direction: column;
    justify-content: center;
  }
  h1 {
    color: #0f172a;
    font-size: 58px;
  }
  h2 {
    color: #0f172a;
    font-size: 44px;
  }
  h3 {
    color: #2563eb;
    font-size: 32px;
  }
  strong {
    color: #2563eb;
  }
  code {
    font-family: "Cascadia Code", monospace;
  }
  table {
    font-size: 24px;
  }
  blockquote {
    border-left: 8px solid #2563eb;
    color: #1f2937;
    font-size: 34px;
    padding-left: 28px;
  }
  section.split h1 {
    margin-bottom: 28px;
  }
  section.split .cols {
    align-items: stretch;
    display: grid;
    gap: 36px;
    grid-template-columns: 0.95fr 1.05fr;
  }
  section.split .cols.equal {
    grid-template-columns: 1fr 1fr;
  }
  section.split .col {
    min-width: 0;
  }
  section.split .panel,
  section.split .code-card,
  section.split .diagram {
    background: #f8fafc;
    border: 1px solid #cbd5e1;
    border-radius: 8px;
    padding: 22px;
  }
  section.split .label {
    color: #2563eb;
    font-size: 22px;
    font-weight: 700;
    margin-bottom: 14px;
    text-transform: uppercase;
  }
  section.split .caption {
    color: #475569;
    font-size: 24px;
    margin-top: 16px;
  }
  section.split pre {
    background: #0f172a;
    border-radius: 8px;
    color: #e5e7eb;
    font-size: 20px;
    line-height: 1.35;
    margin: 0;
    padding: 18px;
  }
  section.split .node,
  section.split .hub,
  section.split .metric {
    background: #ffffff;
    border: 1px solid #94a3b8;
    border-radius: 8px;
    padding: 14px 18px;
  }
  section.split .node {
    font-size: 24px;
    font-weight: 700;
    text-align: center;
  }
  section.split .arrow {
    color: #64748b;
    font-size: 28px;
    font-weight: 800;
    margin: 10px 0;
    text-align: center;
  }
  section.split .flow-grid {
    display: grid;
    gap: 12px;
  }
  section.split .hub {
    border-color: #2563eb;
    color: #0f172a;
    font-size: 28px;
    font-weight: 800;
    margin: 16px 0;
    text-align: center;
  }
  section.split .spokes {
    display: grid;
    gap: 12px;
    grid-template-columns: 1fr 1fr;
  }
  section.split .mini {
    color: #334155;
    font-size: 22px;
  }
  section.split .metric {
    align-items: center;
    display: flex;
    flex-direction: column;
    gap: 12px;
    justify-content: center;
    min-height: 280px;
  }
  section.split .metric strong {
    color: #16a34a;
    font-size: 96px;
    line-height: 1;
  }
  section.split .bar {
    align-items: center;
    display: grid;
    gap: 12px;
    grid-template-columns: 110px 1fr 110px;
    margin: 12px 0;
  }
  section.split .bar-fill {
    background: #2563eb;
    border-radius: 8px;
    height: 18px;
  }
  section.split .bar-fill.green {
    background: #16a34a;
  }
---

<!-- _class: lead -->

# TTSC

### TypeScript-Go ToolChain

Samchon

2026-06-26

---

# Preface

![](https://ttsc.dev/og.jpg)

---

# 1.1. TypeScript-Go

대충 TypeScript-Go 곧 출시한다는 이야기.

대충 성능 10배 빠르다는 이야기.

대충 Go 언어 기반이라 생태계 다 뒤집힌다는 이야기.

---

# 1.2. Transformer

대충 트랜스포머가 뭔지 설명

---

# 1.2. Transformer

### Generic function call with `IMember` type

```typescript
import typia, { tags } from "typia";

interface IMember {
  id: string & tags.Format<"uuid">;
  email: string & tags.Format<"email">;
  age: number &
    tags.Type<"uint32"> &
    tags.ExclusiveMinimum<19> &
    tags.Maximum<100>;
}
typia.createIs<IMember>();
```

---

# 1.2. Transformer

### Becomes transformed JS code

```javascript
import * as _b from "typia/lib/internal/_isFormatEmail";
import * as _a from "typia/lib/internal/_isFormatUuid";
import * as _c from "typia/lib/internal/_isTypeUint32";

(() => {
  const _io0 = (input) =>
    "string" === typeof input.id &&
    _a._isFormatUuid(input.id) &&
    "string" === typeof input.email &&
    _b._isFormatEmail(input.email) &&
    "number" === typeof input.age &&
    _c._isTypeUint32(input.age) &&
    19 < input.age &&
    input.age <= 100;
  return (input) => "object" === typeof input && null !== input && _io0(input);
})();
```

---

# 1.2. Transformer

```typescript
import { TypedBody, TypedRoute } from "@nestia/core";
import { Controller } from "@nestjs/common";

@Controller()
export class ShoppingSaleController {
  @TypedRoute.Post()
  public async create(
    @TypedBody() body: IShoppingSale.ICreate
  ): Promise<IShoppingSale> { ... }
}
```

---

# 1.2. Transformer

### Required "ts-patch"

```json
{
  "scripts": {
    "prepare": "ts-patch install"
  },
  "devDependencies": {
    "ts-patch": "^3.2.1"
  }
}
```

---

# 1.2. Transformer

### Required tsconfig.json configuration

```json
{
  "compilerOptions": {
    "strict": true,
    "plugins": [
      { "transform": "typia/lib/transform" },
      { "transform": "@nestia/core/lib/transform" },
      { "transform": "@nestia/sdk/lib/transform" },
    ]
  }
}
```

---

# 1.2. Transformer

대충 TypeScript-Go 에서 통하지 않는다는 이야기

그래서 좆되었다는 이야기

---

# 2. TTSC

대충 좆되서 본인이 직접 만들었단 이야기

---

# 2.1. Transformer

대충 go 언어 기반의 트랜스포머 생태계 다시 만들었다는 이야기

tsconfig.json 설정 안해도 된다는 이야기

`ttsc` 명령어로 컴파일하면 된다는 이야기

---

# 2.2. Runtime

`ttsx src/index.ts` 로 실행할 수 있단 이야기

대충 `ts-node`보다 10배 빠른데 타입 체킹한다는 이야기

---

# 2.3. Linter

대충 Transformer 플러그인의 형태로 개발했다는 이야기

ttsc 명령어 한 방에 컴파일 에러와 린트 규칙 위반이 한번에 잡힌다는 이야기

이론상 비용이 0에 수렴해, 실제 800배 성능차를 보인다는 이야기

---

# 2.4. Graph

클로드 코드나 코덱스에게 컴파일러 AST 정보 기반 그래프 정보를 준다는 이야기

그래서 grep와 과도한 파일 리드를 막아준다는 이야기

토큰 소모량이 100배로 준다는 이야기

---

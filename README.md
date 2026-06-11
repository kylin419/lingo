<!DOCTYPE html>
<html lang="zh-Hant">
<head>
    <meta charset="UTF-8">
    <style>
        body {
            font-family: -apple-system, BlinkMacSystemFont, "Segoe UI", Helvetica, Arial, sans-serif;
            line-height: 1.6;
            color: #24292e;
            max-width: 800px;
            margin: 40px auto;
            padding: 0 20px;
        }
        .header { text-align: center; margin-bottom: 40px; }
        .badge { margin: 5px; vertical-align: middle; }
        pre {
            background-color: #f6f8fa;
            padding: 16px;
            border-radius: 6px;
            overflow: auto;
        }
        code { font-family: monospace; }
        hr { border: 0; border-top: 1px solid #eaecef; margin: 30px 0; }
        h2 { border-bottom: 1px solid #eaecef; padding-bottom: 10px; }
        blockquote { color: #6a737d; border-left: 4px solid #dfe2e5; padding-left: 16px; margin: 0; }
    </style>
</head>
<body>

    <div class="header">
        <h1>Lingo Bot</h1>
        <p><b>跨語言溝通與視覺識別的智慧解決方案</b></p>
        <div>
            <img src="https://img.shields.io/badge/go-1.22+-blue.svg?style=for-the-badge&logo=go" class="badge" alt="Go Version">
            <img src="https://img.shields.io/badge/license-MIT-green.svg?style=for-the-badge" class="badge" alt="License">
            <img src="https://img.shields.io/badge/status-active-brightgreen.svg?style=for-the-badge" class="badge" alt="Status">
        </div>
    </div>

    <hr>

    <p><strong>Lingo Bot</strong> 是一款專為高效溝通設計的智慧機器人，深度整合頂尖翻譯引擎與 OCR 影像辨識技術，協助您打破語言藩籬，實現即時的文件數位化與跨國團隊協作。</p>

    <h2>核心特色</h2>
    <ul>
        <li><strong>即時多語翻譯</strong>：精準辨識原文語言，並同步翻譯至多種目標語言。</li>
        <li><strong>高精準 OCR 識別</strong>：自動提取圖片中的文字內容，並進行即時翻譯與轉譯。</li>
        <li><strong>直覺化指令交互</strong>：簡潔的 <code>$set</code> 指令系統，快速切換翻譯偏好，適配多種使用情境。</li>
        <li><strong>高效能架構</strong>：為群組環境最佳化，確保大規模併發下的低延遲回應。</li>
    </ul>

    <h2>使用指南</h2>
    <h3>1. 初始化設定</h3>
    <p>設定您的翻譯目標語言（支援 ISO 639-1 標準碼）：</p>
    <pre><code>$set [語言1] [語言2] [語言3] ...</code></pre>
    <blockquote>範例：輸入 <code>$set zh-Hant en ja</code>，系統將自動把訊息翻譯為繁體中文、英文與日文。</blockquote>

    <h3>2. 即時翻譯服務</h3>
    <p>完成設定後，輸入文字訊息，機器人將自動觸發翻譯邏輯，呈現對應結果。</p>

    <h3>3. 圖像文字處理</h3>
    <p>直接上傳圖片，Lingo Bot 將自動識別圖中文字，並根據您的語言設定進行翻譯回傳。</p>

    <hr>

    <h2>聯繫與協作</h2>
    <p>本專案由 <strong>Kylin</strong> 持續維護。歡迎透過 GitHub Issues 分享您的使用反饋。</p>

    <div style="text-align: center; color: #6a737d; margin-top: 50px;">
        <small>Powered by Azure & LINE API</small>
    </div>

</body>
</html>

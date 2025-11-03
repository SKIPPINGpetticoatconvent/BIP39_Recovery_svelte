<script lang="ts">
  import {InitBIP39, GetWord, IsValidNumber, CalculateWordIndex, ValidateWordIndex} from '../wailsjs/go/main/App.js'
  import {tick} from 'svelte'

  // Application state
  let currentLang = 'zh'
  let mnemonicLength = 12
  let currentWordIndex = 0
  let recoveredWords: string[] = []
  let currentWordSum = 0
  let currentWordInputs: number[] = []
  let inputNumber = ''

  // UI state - Page navigation
  let currentPage = 'welcome' // 'welcome', 'recovery', 'result'

  // Languages
  const languages = {
    en: {
      window_title: "Offline BIP39 Mnemonic Recovery Tool",
      welcome_header: "BIP39 Mnemonic Recovery",
      select_length_prompt: "Please select the length of your seed phrase:",
      "12_words": "12 Words",
      "18_words": "18 Words",
      "24_words": "24 Words",
      offline_warning: "This tool is 100% offline. No data is ever sent.",
      wordlist_file_error_title: "Wordlist File Error",
      wordlist_not_found: "Wordlist file '{filename}' not found!\n\nPlease ensure it is in the same directory as the script.",
      wordlist_invalid_length: "The wordlist '{filename}' is invalid.\n\nIt contains {count} words, but it must contain exactly 2048.",
      file_read_error_title: "File Read Error",
      file_read_error_message: "An error occurred while reading the file: {error}",
      recovering_word_title: "Recovering Word {current} of {total}",
      enter_number_label: "Enter number (e.g., 2, 4, 256):",
      add_number_button: "Add Number",
      entered_numbers_label: "Entered Numbers: {numbers}",
      current_word_label: "Current Word: {status}",
      status_waiting: "(waiting for input)",
      status_invalid_index: "[Sum: {sum}] -> INVALID INDEX",
      status_valid_word: "[Sum: {sum}] -> Index {index} -> '{word}'",
      confirm_and_next_button: "Confirm Word & Next",
      recovered_words_header: "Recovered Words so far:",
      invalid_input_title: "Invalid Input",
      invalid_input_int_warning: "Please enter a valid whole number.",
      invalid_input_power_of_2_warning: "Please enter a valid power of 2 (1, 2, 4, ..., 1024).",
      duplicate_input_warning: "The number {num} has already been added for this word.",
      no_input_title: "No Input",
      no_input_warning: "Please add at least one number for this word.",
      sum_error_title: "Error",
      sum_error_message: "The sum of the numbers is invalid and does not correspond to a valid word.",
      recovery_complete_header: "Recovery Successful!",
      your_seed_phrase_is: "Your recovered BIP39 seed phrase is:",
      security_note: "SECURITY NOTE: Please close this window after you have secured your phrase.",
      restart_button: "Restart",
      quit_button: "Quit",
    },
    zh: {
      window_title: "离线BIP39助记词恢复工具",
      welcome_header: "BIP39 助记词恢复",
      select_length_prompt: "请选择您的助记词短语长度：",
      "12_words": "12个单词",
      "18_words": "18个单词",
      "24_words": "24个单词",
      offline_warning: "本工具为100%离线工具，绝不发送任何数据。",
      wordlist_file_error_title: "词库文件错误",
      wordlist_not_found: "词库文件 '{filename}' 未找到！\n\n请确保该文件与脚本在同一目录下。",
      wordlist_invalid_length: "词库文件 '{filename}' 无效。\n\n它包含 {count} 个单词，但必须是2048个。",
      file_read_error_title: "文件读取错误",
      file_read_error_message: "读取文件时发生错误: {error}",
      recovering_word_title: "正在恢复第 {current} / {total} 个单词",
      enter_number_label: "输入数字 (例如 2, 4, 256):",
      add_number_button: "添加数字",
      entered_numbers_label: "已输入的数字: {numbers}",
      current_word_label: "当前单词: {status}",
      status_waiting: "(等待输入)",
      status_invalid_index: "[总和: {sum}] -> 无效索引",
      status_valid_word: "[总和: {sum}] -> 索引 {index} -> '{word}'",
      confirm_and_next_button: "确认单词并继续",
      recovered_words_header: "已恢复的单词:",
      invalid_input_title: "无效输入",
      invalid_input_int_warning: "请输入一个有效的整数。",
      invalid_input_power_of_2_warning: "请输入一个有效的2的幂（1, 2, 4, ..., 1024）。",
      duplicate_input_warning: "数字 {num} 已经为这个单词添加过了。",
      no_input_title: "没有输入",
      no_input_warning: "请至少为这个单词添加一个数字。",
      sum_error_title: "错误",
      sum_error_message: "数字总和无效，无法对应到一个有效的单词。",
      recovery_complete_header: "恢复成功！",
      your_seed_phrase_is: "您恢复的BIP39助记词短语是：",
      security_note: "【安全提示】在您安全备份好助记词后，请关闭本窗口。",
      restart_button: "重新开始",
      quit_button: "退出",
    },
  }

  // Reactive translation function - updates automatically when currentLang changes
  $: T = (key: string): string => {
    return languages[currentLang as keyof typeof languages]?.[key as keyof (typeof languages)[keyof typeof languages]] as string || key
  }

  // Initialize BIP39 on startup
  async function initialize() {
    try {
      await InitBIP39()
    } catch (error) {
      console.error('Failed to initialize BIP39:', error)
      alert(T('wordlist_file_error_title') + ': ' + T('wordlist_invalid_length').replace('{count}', 'unknown'))
    }
  }

  initialize()

  function setLanguage(lang: string) {
    currentLang = lang
    // Svelte will automatically re-render components that use currentLang
  }

  function startRecovery(length: number) {
    mnemonicLength = length
    currentWordIndex = 0
    recoveredWords = [] // Ensure reactivity
    resetCurrentWord()
    currentPage = 'recovery'
  }

  function resetCurrentWord() {
    currentWordSum = 0
    currentWordInputs = [] // Ensure reactivity
    inputNumber = ''
  }

  async function addNumber() {
    const numStr = inputNumber.trim()
    if (!numStr) return

    try {
      const num = parseInt(numStr)
      if (!await IsValidNumber(num)) {
        alert(T('invalid_input_power_of_2_warning'))
        return
      }

      if (currentWordInputs.includes(num) && num !== 1) {
        alert(T('duplicate_input_warning').replace('{num}', num.toString()))
        return
      }

      if (currentWordInputs.filter(n => n === 1).length >= 3) {
        alert("数字1最多只能重复两次。")
        return
      }

      currentWordInputs.push(num)
      currentWordInputs.sort((a, b) => a - b)
      currentWordSum += num
      inputNumber = ''
    } catch (error) {
      alert(T('invalid_input_int_warning'))
    }
  }

  async function processNextWord() {
    if (currentWordInputs.length === 0) {
      alert(T('no_input_warning'))
      return
    }

    const wordIndex = await CalculateWordIndex(currentWordSum)
    if (!await ValidateWordIndex(wordIndex)) {
      alert(T('sum_error_message'))
      return
    }

    try {
      const word = await GetWord(wordIndex)
      recoveredWords = [...recoveredWords, word] // Trigger reactivity
      currentWordIndex++

      // Update status with actual word
      currentWordStatus = T('status_valid_word')
        .replace('{sum}', currentWordSum.toString())
        .replace('{index}', (wordIndex + 1).toString())
        .replace('{word}', word)

      if (currentWordIndex >= mnemonicLength) {
        showFinalResult()
      } else {
        resetCurrentWord()
      }
    } catch (error) {
      console.error('Error getting word:', error)
      alert(T('sum_error_message'))
    }
  }

  function showFinalResult() {
    currentPage = 'result'
  }

  function restart() {
    currentPage = 'welcome'
    mnemonicLength = 12
    currentWordIndex = 0
    recoveredWords = [] // Ensure reactivity
    resetCurrentWord()
  }

  async function getCurrentWordStatus() {
    if (currentWordSum === 0) {
      return T('status_waiting')
    }
    const wordIndex = currentWordSum - 1
    if (wordIndex >= 0 && wordIndex < 2048) {
      try {
        const word = await GetWord(wordIndex)
        return T('status_valid_word')
          .replace('{sum}', currentWordSum.toString())
          .replace('{index}', (wordIndex + 1).toString())
          .replace('{word}', word)
      } catch (error) {
        return T('status_invalid_index').replace('{sum}', currentWordSum.toString())
      }
    } else {
      return T('status_invalid_index').replace('{sum}', currentWordSum.toString())
    }
  }

  // Reactive status
  let currentWordStatus = ''
  $: {
    currentWordStatus = T('status_waiting')
    if (currentWordSum > 0) {
      const wordIndex = currentWordSum - 1
      if (wordIndex >= 0 && wordIndex < 2048) {
        // Show index info without waiting for word retrieval
        currentWordStatus = T('status_valid_word')
          .replace('{sum}', currentWordSum.toString())
          .replace('{index}', (wordIndex + 1).toString())
          .replace('{word}', T('status_waiting'))
      } else {
        currentWordStatus = T('status_invalid_index').replace('{sum}', currentWordSum.toString())
      }
    }
  }

  function quit() {
    window.close()
  }
</script>

<main class="app">
  <!-- Language Switcher -->
  {#if currentPage !== ''}
    <div class="language-switcher">
      <button class="lang-btn {currentLang === 'en' ? 'active' : ''}" on:click={() => setLanguage('en')}>English</button>
      <button class="lang-btn {currentLang === 'zh' ? 'active' : ''}" on:click={() => setLanguage('zh')}>中文</button>
    </div>
  {/if}

  <!-- Welcome Screen -->
  {#if currentPage === 'welcome'}
    <div class="screen welcome">
      <h1>{T('welcome_header')}</h1>
      <p>{T('select_length_prompt')}</p>
      <div class="length-buttons">
        <button class="length-btn" on:click={() => startRecovery(12)}>{T('12_words')}</button>
        <button class="length-btn" on:click={() => startRecovery(18)}>{T('18_words')}</button>
        <button class="length-btn" on:click={() => startRecovery(24)}>{T('24_words')}</button>
      </div>
      <p class="offline-warning">{T('offline_warning')}</p>
    </div>
  {/if}

  <!-- Recovery Screen -->
  {#if currentPage === 'recovery'}
    <div class="screen recovery">
      <h1>{T('recovering_word_title').replace('{current}', (currentWordIndex + 1).toString()).replace('{total}', mnemonicLength.toString())}</h1>

      <div class="input-section">
        <label for="number-input">{T('enter_number_label')}</label>
        <div class="input-row">
          <input
            id="number-input"
            type="text"
            bind:value={inputNumber}
            on:keydown={(e) => e.key === 'Enter' && addNumber()}
            placeholder="2, 4, 256..."
          />
          <button class="add-btn" on:click={addNumber}>{T('add_number_button')}</button>
        </div>
      </div>

      <div class="status-section">
        <p>{T('entered_numbers_label').replace('{numbers}', currentWordInputs.join(', '))}</p>
        <p>{T('current_word_label').replace('{status}', currentWordStatus)}</p>
      </div>

      <button class="confirm-btn" on:click={processNextWord}>{T('confirm_and_next_button')}</button>

      <div class="words-section">
        <h3>{T('recovered_words_header')}</h3>
        <div class="words-list">
          {#each recoveredWords as word, i}
            <div class="word-item">
              <span class="word-number">{i + 1}.</span>
              <span class="word-text">{word}</span>
            </div>
          {/each}
        </div>
      </div>
    </div>
  {/if}

  <!-- Result Screen -->
  {#if currentPage === 'result'}
    <div class="screen result">
      <h1>{T('recovery_complete_header')}</h1>
      <p>{T('your_seed_phrase_is')}</p>
      <div class="final-words">
        {#each recoveredWords as word, i}
          <div class="final-word-item">
            <span class="word-number">{i + 1}.</span>
            <span class="word-text">{word}</span>
          </div>
        {/each}
      </div>
      <p class="security-note">{T('security_note')}</p>
      <div class="result-buttons">
        <button class="restart-btn" on:click={restart}>{T('restart_button')}</button>
        <button class="quit-btn" on:click={quit}>{T('quit_button')}</button>
      </div>
    </div>
  {/if}
</main>

<style>
  .app {
    font-family: 'Segoe UI', Tahoma, Geneva, Verdana, sans-serif;
    background: #f5f5f5;
    min-height: 100vh;
    padding: 20px;
    color: #333;
    display: flex;
    flex-direction: column;
    align-items: center;
  }

  .language-switcher {
    position: fixed;
    top: 20px;
    right: 20px;
    display: flex;
    gap: 10px;
    z-index: 1000;
  }

  .lang-btn {
    padding: 8px 16px;
    border: 1px solid #ddd;
    background: rgba(255, 255, 255, 0.9);
    border-radius: 6px;
    cursor: pointer;
    transition: all 0.3s ease;
  }

  .lang-btn.active {
    background: #4CAF50;
    color: white;
    border-color: #4CAF50;
  }

  .screen {
    max-width: 550px;
    width: 100%;
    margin: 80px auto 0;
    background: #ffffff;
    padding: 40px;
    border-radius: 8px;
    box-shadow: 0 4px 8px rgba(0, 0, 0, 0.1);
    text-align: center;
  }

  h1 {
    color: #2c3e50;
    margin-bottom: 20px;
    font-size: 2rem;
  }

  .welcome p {
    margin: 20px 0;
    color: #666;
    font-size: 1.1rem;
  }

  .length-buttons {
    display: flex;
    gap: 15px;
    justify-content: center;
    margin: 30px 0;
  }

  .length-btn {
    padding: 12px 24px;
    background: #ffffff;
    color: #333;
    border: 1px solid #dee2e6;
    border-radius: 6px;
    cursor: pointer;
    font-size: 1rem;
    transition: all 0.2s ease;
  }

  .length-btn:hover {
    background: #f8f9fa;
    border-color: #007bff;
  }

  .length-btn:hover {
    transform: translateY(-2px);
  }

  .offline-warning {
    color: #e74c3c;
    font-weight: bold;
    margin-top: 30px;
  }

  .input-section {
    margin: 30px 0;
  }

  .input-section label {
    display: block;
    margin-bottom: 10px;
    font-weight: bold;
    color: #2c3e50;
  }

  .input-row {
    display: flex;
    gap: 10px;
    align-items: center;
  }

  .input-row input {
    flex: 1;
    padding: 12px;
    border: 2px solid #ddd;
    border-radius: 6px;
    font-size: 1rem;
  }

  .input-row input:focus {
    outline: none;
    border-color: #667eea;
  }

  .add-btn {
    padding: 12px 20px;
    background: #4CAF50;
    color: white;
    border: none;
    border-radius: 6px;
    cursor: pointer;
    font-weight: bold;
  }

  .add-btn:hover {
    background: #45a049;
  }

  .status-section {
    margin: 20px 0;
    padding: 15px;
    background: #f8f9fa;
    border-radius: 6px;
  }

  .status-section p {
    margin: 5px 0;
    font-family: 'Courier New', monospace;
    font-weight: bold;
  }

  .confirm-btn {
    padding: 14px 28px;
    background: #007bff;
    color: white;
    border: none;
    border-radius: 6px;
    cursor: pointer;
    font-size: 1.1rem;
    font-weight: bold;
    margin: 20px 0;
    transition: all 0.2s ease;
  }

  .confirm-btn:hover {
    background: #0056b3;
  }

  .words-section h3 {
    color: #2c3e50;
    margin: 30px 0 15px;
  }

  .words-list {
    background: #f8f9fa;
    border-radius: 8px;
    padding: 15px;
    max-height: 200px;
    overflow-y: auto;
  }

  .word-item {
    display: flex;
    align-items: center;
    margin: 8px 0;
    padding: 8px 12px;
    background: white;
    border-radius: 4px;
    box-shadow: 0 1px 3px rgba(0, 0, 0, 0.1);
  }

  .word-number {
    font-weight: bold;
    margin-right: 10px;
    color: #667eea;
    min-width: 25px;
  }

  .word-text {
    font-family: 'Courier New', monospace;
    color: #2c3e50;
  }

  .final-words {
    display: grid;
    grid-template-columns: repeat(auto-fit, minmax(200px, 1fr));
    gap: 10px;
    margin: 30px 0;
    padding: 20px;
    background: #f8f9fa;
    border-radius: 8px;
  }

  .final-word-item {
    display: flex;
    align-items: center;
    padding: 10px 15px;
    background: white;
    border-radius: 6px;
    box-shadow: 0 2px 4px rgba(0, 0, 0, 0.1);
    font-family: 'Courier New', monospace;
  }

  .final-word-item .word-number {
    color: #4CAF50;
    font-weight: bold;
  }

  .final-word-item .word-text {
    color: #2c3e50;
    font-weight: bold;
  }

  .security-note {
    color: #e74c3c;
    font-weight: bold;
    margin: 20px 0;
    font-size: 0.9rem;
  }

  .result-buttons {
    display: flex;
    gap: 15px;
    justify-content: center;
    margin-top: 30px;
  }

  .restart-btn {
    padding: 12px 24px;
    background: #4CAF50;
    color: white;
    border: none;
    border-radius: 6px;
    cursor: pointer;
    font-weight: bold;
  }

  .quit-btn {
    padding: 12px 24px;
    background: #e74c3c;
    color: white;
    border: none;
    border-radius: 6px;
    cursor: pointer;
    font-weight: bold;
  }

  .restart-btn:hover {
    background: #45a049;
  }

  .quit-btn:hover {
    background: #c0392b;
  }
</style>

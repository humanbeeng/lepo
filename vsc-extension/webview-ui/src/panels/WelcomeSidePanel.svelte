<script lang="ts">
  let name: string = " @humanbeeng";
</script>

<main class="flex h-screen">
  <div class="flex flex-col h-screen justify-between">
    <div class="justify-start">
      <div class="py-5">
        <img src="webview-ui/src/media/smile.svg" alt="Lepo Logo" />
        <h3 class="text-base font-bold">Lepo Chat</h3>
      </div>
      <p class="text-sm">
        Hello
        <span class="font-extrabold">{name}</span> !<br /> I'm <span class="font-semibold">Lepo</span>, your
        AI assistant who can help you with queries regarding your codebase. I can explain given code, help you
        find and build things using your code's context.

        <br />
        <br />

        You can ask me any question about this codebase and I'll try my best to answer. <br /><br />
        Please note that since I'm an AI assistant, I may not have all knowledge about your project. Make sure
        to verify my answers and iterate upon that.
      </p>
    </div>
    <div class="justify-end align-bottom py-3">
      <div class="flex w-full justify-between gap-2 flex-row">
        <vscode-text-area />
        <vscode-button>Send</vscode-button>
      </div>
    </div>
  </div>
</main>

<script lang="ts">
  let name: string = " @humanbeeng";
  let text: string = "";

  const handleSubmit = () => {
    console.log(text);
  };
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
      <p>{text}</p>
    </div>
    <div class="py-3 bg-green-500">
      <div class="flex w-full justify-between gap-2">
        <vscode-text-area
          value={text}
          on:input={(event) => {
            text = event.target.value;
          }}
          class="w-full h-full"
        />
        <vscode-button on:keypress={handleSubmit} on:click={handleSubmit}>Send</vscode-button>
      </div>
    </div>
  </div>
</main>

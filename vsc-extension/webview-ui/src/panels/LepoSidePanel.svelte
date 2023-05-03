<script lang="ts">
  import { vscode } from "../utilities/vscode";

  export let onclick: () => void;

  const login = () => {
    vscode.postMessage({ command: "hello", text: "Calling from sidepanel" });
    onclick();
  };
</script>

<main>
  <h1>Hello from sidepanel</h1>
  <vscode-button on:click={login} on:keypress={login}>Hello</vscode-button>
</main>

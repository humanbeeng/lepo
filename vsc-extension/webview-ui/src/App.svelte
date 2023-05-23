<script lang="ts">
  import {
    provideVSCodeDesignSystem,
    vsCodeButton,
    vsCodeTextArea,
    vsCodeTextField,
  } from "@vscode/webview-ui-toolkit";
  import ChatSidePanel from "./panels/ChatSidePanel.svelte";
  import LepoSidePanel from "./panels/LepoSidePanel.svelte";
  import { vscode } from "./utilities/vscode";

  provideVSCodeDesignSystem().register(vsCodeButton(), vsCodeTextField(), vsCodeTextArea());

  let page: "welcome" | "chat" = vscode.getState()?.page || "welcome";

  $: {
    vscode.setState({ page });
  }

  export const handleHello = () => {
    console.log("Handling hello");
    page = "chat";
  };
</script>

<main>
  {#if page === "chat"}
    <LepoSidePanel
      onclick={() => {
        page = "welcome";
      }}
    />
  {:else if page === "welcome"}
    <ChatSidePanel />
  {/if}
</main>

<style lang="postcss" global>
  @tailwind base;
  @tailwind components;
  @tailwind utilities;
</style>

import * as vscode from "vscode";
import { CancellationToken, WebviewView, WebviewViewProvider, WebviewViewResolveContext } from "vscode";
import { Message } from "../../model/Message";
import { getNonce } from "../../utilities/getNonce";
import { getUri } from "../../utilities/getUri";

export default class SidePanelProvider implements WebviewViewProvider {
  private webView?: vscode.WebviewView;

  constructor(private context: vscode.ExtensionContext, private readonly extensionUri: vscode.Uri) {}

  resolveWebviewView(
    webviewView: WebviewView,
    context: WebviewViewResolveContext<unknown>,
    token: CancellationToken
  ): void | Thenable<void> {
    this.webView = webviewView;
    webviewView.webview.html = this.getHTMLForWebview(webviewView.webview);
    webviewView.webview.options = {
      enableScripts: true,
      localResourceRoots: [this.context.extensionUri],
    };
    this._setWebviewMessageListener(webviewView);
  }

  private _setWebviewMessageListener(webviewView: WebviewView) {
    webviewView.webview.onDidReceiveMessage((message: Message) => {
      const command = message.command;
      const text = message.text;
      switch (command) {
        case "hello":
          vscode.window.showInformationMessage("Hello clicked");
      }
    });
  }

  private getHTMLForWebview(webview: vscode.Webview) {
    // The CSS file from the Svelte build output
    const stylesUri = getUri(webview, this.extensionUri, ["webview-ui", "public", "build", "bundle.css"]);
    // The JS file from the Svelte build output
    const scriptUri = getUri(webview, this.extensionUri, ["webview-ui", "public", "build", "bundle.js"]);

    const nonce = getNonce();

    return `<!DOCTYPE html>
	  <html lang="en">
		<head>
		  <meta charset="UTF-8" />
		  <meta http-equiv="X-UA-Compatible" content="IE=edge" />
		  <meta http-equiv="Content-Security-Policy" content="default-src  img-src https: data:; style-src 'unsafe-inline' ${webview.cspSource}; script-src 'nonce-${nonce}';">
					<meta name="viewport" content="width=device-width, initial-scale=1.0">
			<link href="${stylesUri}" rel="stylesheet"/>
		</head>
	  
	  <body><script nonce="${nonce}" src="${scriptUri}"></script></body>
	  </html>`;
  }
}

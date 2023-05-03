import * as vscode from "vscode";
import { ExtensionContext } from "vscode";
import LepoSidePanelProvider from "./provider/sidepanel/LepoSidePanelProvider";

export async function activate(context: ExtensionContext) {
  const lepoSidePanelProvider = new LepoSidePanelProvider(context, context.extensionUri);
  vscode.window.registerWebviewViewProvider("lepo.main", lepoSidePanelProvider, {});
}

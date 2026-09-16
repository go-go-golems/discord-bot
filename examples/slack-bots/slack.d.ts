declare module "slack" {
  export type JSONValue = null | boolean | number | string | JSONValue[] | {[key: string]: JSONValue};
  export interface MessageRef {channelId: string; ts: string}
  export interface Text {text: string}
  export interface Block {type: string; [key: string]: JSONValue}
  export interface MessagePayload {text: string; blocks?: Block[]}
  export interface PostMessage extends MessagePayload {channelId: string; threadTs?: string}
  export interface Context {
    id: string; teamId: string; channelId: string; userId: string;
    command: string; text: string;
    event: {type: string; text: string; ts: string; threadTs: string; channelId: string; userId: string};
    action?: {type: string; actionId: string; blockId?: string; value?: string; selectedOption?: Record<string, JSONValue>; selectedOptions?: JSONValue[]; selection?: Record<string, JSONValue>; messageTs?: string; threadTs?: string};
    view?: {type: string; callbackId: string; privateMetadata?: string; id?: string; hash?: string};
    values?: {all: Record<string, Record<string, JSONValue>>; text(blockId: string, actionId: string): string | undefined};
    shortcut?: {type: string; callbackId: string; message?: Record<string, JSONValue>};
    query?: string;
    replaceOriginal(message: MessagePayload): Promise<void>;
    config: Record<string, string | boolean | number>;
    reply(message: MessagePayload): Promise<MessageRef | {delivered: true; via: "response_url"}>;
    slack: {
      messages: {post(message: PostMessage): Promise<MessageRef>; update(params: Record<string, JSONValue>): Promise<Record<string, JSONValue>>; delete(params: Record<string, JSONValue>): Promise<Record<string, JSONValue>>; ephemeral(params: Record<string, JSONValue>): Promise<Record<string, JSONValue>>; permalink(params: Record<string, JSONValue>): Promise<Record<string, JSONValue>>};
      conversations: Record<"history" | "replies" | "info" | "list" | "members" | "join" | "leave" | "setTopic" | "kick" | "archive", (params: Record<string, JSONValue>) => Promise<Record<string, JSONValue>>>;
      users: Record<"info" | "list", (params: Record<string, JSONValue>) => Promise<Record<string, JSONValue>>>;
      pins: Record<"add" | "remove" | "list", (params: Record<string, JSONValue>) => Promise<Record<string, JSONValue>>>;
      reactions: Record<"add" | "remove" | "get", (params: Record<string, JSONValue>) => Promise<Record<string, JSONValue>>>;
      usergroups: Record<"list" | "members" | "setMembers", (params: Record<string, JSONValue>) => Promise<Record<string, JSONValue>>>;
      workspace: {info(params: Record<string, JSONValue>): Promise<Record<string, JSONValue>>};
      files: {upload(params: {channel_id: string; filename: string; content: string; thread_ts?: string}): Promise<Record<string, JSONValue>>};
    };
    store: {get(key: string): JSONValue | undefined; set(key: string, value: JSONValue): void; delete(key: string): boolean; keys(): string[]};
    log: {debug(message: string): void; info(message: string): void; warn(message: string): void; error(message: string): void};
    ack?: {accept(): Promise<void>; errors(errors: Record<string, string>): Promise<void>; update(view: Block): Promise<void>; options(options: Block[]): Promise<void>};
    openModal(view: Block): Promise<{id: string; hash: string}>;
  }
  export interface Registration {
    configure(spec: {name: string; description?: string; scopes?: string[]; run?: {fields: Record<string, {type: "string" | "bool" | "number"; default?: string | boolean | number; required?: boolean; help?: string}>}}): void;
    command(name: string, spec: {description: string}, handler: (ctx: Context) => MessagePayload | void | Promise<MessagePayload | void>): void;
    shortcut(spec: {callbackId: string; name: string; description: string; type: "message" | "global"}, handler: (ctx: Context) => void | Promise<void>): void;
    options(actionId: string, handler: (ctx: Context) => void | Promise<void>): void;
    event(name: "app_mention" | "message" | "reaction_added" | "reaction_removed" | "member_joined_channel" | "member_left_channel" | "team_join", handler: (ctx: Context) => MessagePayload | void | Promise<MessagePayload | void>): void;
    view(callbackId: string, handler: (ctx: Context) => void | Promise<void>): void;
    action(actionId: string, handler: (ctx: Context) => MessagePayload | void | Promise<MessagePayload | void>): void;
  }
  export function defineBot(register: (api: Registration) => void): object;
}

declare module "slack/ui" {
  import {Block} from "slack";
  export function plain(text: string): Block;
  export function mrkdwn(text: string): Block;
  export function divider(): Block;
  export function header(text: string): Block;
  export interface ButtonBuilder {
    value(value: string): ButtonBuilder;
    style(style: "primary" | "danger" | ""): ButtonBuilder;
    build(): Block;
  }
  export function button(actionId: string, label: string): ButtonBuilder;
  export function section(text: Block, accessory?: Block | ButtonBuilder): Block;
  export function actions(blockId: string, ...buttons: (ButtonBuilder | Block)[]): Block;
  export interface MessageBuilder {
    block(block: Block): MessageBuilder;
    build(): {text: string; blocks: Block[]};
  }
  export function message(text: string): MessageBuilder;
  export interface TextInputBuilder {
    initial(value: string): TextInputBuilder;
    placeholder(value: string): TextInputBuilder;
    multiline(): TextInputBuilder;
    length(min: number, max: number): TextInputBuilder;
    build(): Block;
  }
  export interface InputOptions {optional?: boolean; hint?: string; dispatch_action?: boolean}
  export function option(label: string, value: string): Block;
  export function context(...elements: Block[]): Block;
  export function image(url: string, altText: string): Block;
  export function linkButton(actionId: string, label: string, url: string): Block;
  export function confirm(element: Block | ButtonBuilder, title: string, text: string, accept?: string, cancel?: string): Block;
  export function input(blockId: string, label: string, input: TextInputBuilder | Block, options?: InputOptions): Block;
  export function textInput(actionId: string): TextInputBuilder;
  export interface ModalBuilder {
    block(block: Block): ModalBuilder;
    metadata(value: string): ModalBuilder;
    input(blockId: string, label: string, input: TextInputBuilder | Block, options?: InputOptions): ModalBuilder;
    input(block: Block): ModalBuilder;
    submit(label: string): ModalBuilder;
    close(label: string): ModalBuilder;
    build(): Block;
  }
  export function modal(callbackId: string, title: string): ModalBuilder;
  export function staticSelect(actionId: string, options?: Record<string, unknown>): Block;
  export function multiStaticSelect(actionId: string, options?: Record<string, unknown>): Block;
  export function externalSelect(actionId: string, options?: Record<string, unknown>): Block;
  export function multiExternalSelect(actionId: string, options?: Record<string, unknown>): Block;
  export function usersSelect(actionId: string, options?: Record<string, unknown>): Block;
  export function multiUsersSelect(actionId: string, options?: Record<string, unknown>): Block;
  export function channelsSelect(actionId: string, options?: Record<string, unknown>): Block;
  export function multiChannelsSelect(actionId: string, options?: Record<string, unknown>): Block;
  export function conversationsSelect(actionId: string, options?: Record<string, unknown>): Block;
  export function multiConversationsSelect(actionId: string, options?: Record<string, unknown>): Block;
  export function datePicker(actionId: string, options?: Record<string, unknown>): Block;
  export function timePicker(actionId: string, options?: Record<string, unknown>): Block;
  export function checkboxes(actionId: string, options?: Record<string, unknown>): Block;
  export function radioButtons(actionId: string, options?: Record<string, unknown>): Block;
  export function overflow(actionId: string, options?: Record<string, unknown>): Block;
}

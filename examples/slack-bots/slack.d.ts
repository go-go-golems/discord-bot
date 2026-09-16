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
    action?: {type: string; actionId: string; blockId?: string; value?: string; selectedOption?: Record<string, JSONValue>; selectedOptions?: JSONValue[]; messageTs?: string; threadTs?: string};
    view?: {type: string; callbackId: string; privateMetadata?: string; id?: string; hash?: string};
    values?: {all: Record<string, Record<string, JSONValue>>; text(blockId: string, actionId: string): string | undefined};
    config: Record<string, string | boolean | number>;
    reply(message: MessagePayload): Promise<MessageRef | {delivered: true; via: "response_url"}>;
    slack: {messages: {post(message: PostMessage): Promise<MessageRef>}};
    store: {get(key: string): JSONValue | undefined; set(key: string, value: JSONValue): void; delete(key: string): boolean; keys(): string[]};
    log: {debug(message: string): void; info(message: string): void; warn(message: string): void; error(message: string): void};
    ack?: {accept(): Promise<void>; errors(errors: Record<string, string>): Promise<void>};
    openModal(view: Block): Promise<{id: string; hash: string}>;
  }
  export interface Registration {
    configure(spec: {name: string; description?: string; run?: {fields: Record<string, {type: "string" | "bool" | "number"; default?: string | boolean | number; required?: boolean; help?: string}>}}): void;
    command(name: string, spec: {description: string}, handler: (ctx: Context) => MessagePayload | void | Promise<MessagePayload | void>): void;
    event(name: "app_mention", handler: (ctx: Context) => MessagePayload | void | Promise<MessagePayload | void>): void;
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
    style(style: "primary" | "danger"): ButtonBuilder;
    build(): Block;
  }
  export function button(actionId: string, label: string): ButtonBuilder;
  export function section(text: Block): Block;
  export function actions(blockId: string, ...buttons: (ButtonBuilder | Block)[]): Block;
  export interface MessageBuilder {
    block(block: Block): MessageBuilder;
    build(): {text: string; blocks: Block[]};
  }
  export function message(text: string): MessageBuilder;
  export interface TextInputBuilder {
    initial(value: string): TextInputBuilder;
    placeholder(value: string): TextInputBuilder;
    required(): TextInputBuilder;
    optional(): TextInputBuilder;
    build(): Block;
  }
  export function input(blockId: string, label: string, input: TextInputBuilder | Block): Block;
  export function textInput(actionId: string): TextInputBuilder;
  export interface ModalBuilder {
    metadata(value: string): ModalBuilder;
    input(blockId: string, label: string, input: TextInputBuilder | Block): ModalBuilder;
    input(block: Block): ModalBuilder;
    submit(label: string): ModalBuilder;
    close(label: string): ModalBuilder;
    build(): Block;
  }
  export function modal(callbackId: string, title: string): ModalBuilder;
}

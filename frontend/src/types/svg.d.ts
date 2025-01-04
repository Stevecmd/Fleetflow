// Declare module for SVG imports in TypeScript
declare module '*.svg' {
  const content: string;
  export default content;
}

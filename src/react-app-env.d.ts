/// <reference types="react-scripts" />

declare module "*.scss" {
  const content: { [className: string]: string };
  export default content;
}

declare module "*.css";
declare module "*.png";
declare module "*.jpg";
declare module "*.jpeg";
declare module "*.svg";
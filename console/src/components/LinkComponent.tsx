import React from "react";
import { Link, LinkProps } from "react-router-dom";

const LinkComponent = React.forwardRef<
  HTMLAnchorElement,
  Omit<LinkProps, "to"> & { href: LinkProps["to"] }
>((props, ref) => {
  const { href, ...other } = props;
  return <Link to={href} {...other} ref={ref} />;
});

export default LinkComponent;

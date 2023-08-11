import { useState, useCallback } from "react";
import { FallbackProps } from "react-error-boundary";

import Snackbar from "./Snackbar";

export default function ErrorFallback({ error }: FallbackProps) {
  const [open, setOpen] = useState(true);
  const onClose = useCallback(() => setOpen(false), [setOpen]);

  console.log(error);

  return (
    <Snackbar open={open} severity="error" onClose={onClose}>
      {error.message}
    </Snackbar>
  );
}

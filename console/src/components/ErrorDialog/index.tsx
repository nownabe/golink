import { Code, ConnectError } from "@bufbuild/connect";
import { Alert, AlertTitle } from "@mui/material";

type Props = {
  error: Error;
};

export default function ErrorDialog({ error }: Props) {
  const err = ConnectError.from(error);
  console.error(err);

  return (
    <Alert severity="error">
      <AlertTitle>Error: {Code[err.code]}</AlertTitle>
      {err.message}
    </Alert>
  );
}

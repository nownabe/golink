import { useRouteError } from "react-router-dom";

import ErrorDialog from "../ErrorDialog";

export default function RouterErrorDialog() {
  const e = useRouteError() as Error;

  return <ErrorDialog error={e} />;
}

import { Typography } from "@mui/material";
import { Helmet } from "react-helmet";

export default function Home() {
  return (
    <>
      <Helmet>
        <title>Golink Console</title>
      </Helmet>
      <Typography variant="h5" component="h2">
        Welcome to the Golink console!
      </Typography>
    </>
  );
}

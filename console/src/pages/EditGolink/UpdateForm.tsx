import { TextField } from "@mui/material";
import Grid from "@mui/material/Unstable_Grid2";

import ProgressButton from "@/components/ProgressButton";
import Snackbar from "@/components/Snackbar";

import { useUrl } from "./hooks";
import { Golink } from "@/gen/golink/v1/golink_pb";

type Props = {
  golink: Golink;
};

export default function UpdateForm({ golink }: Props) {
  const {
    urlRef,
    urlUpdating,
    onUrlUpdate,
    openUrlUpdateSuccess,
    onUrlUpdateSuccessClose,
    urlUpdateError,
    onUrlUpdateErrorClose,
  } = useUrl(golink.name);

  return (
    <>
      <Grid xs={12}>
        <TextField
          label="URL"
          inputRef={urlRef}
          fullWidth
          defaultValue={golink.url}
        />
      </Grid>
      <Grid xs={12}>
        <ProgressButton loading={urlUpdating} onClick={onUrlUpdate}>
          Update
        </ProgressButton>
      </Grid>
      <Snackbar
        open={!!urlUpdateError}
        severity="error"
        onClose={onUrlUpdateErrorClose}
      >
        {urlUpdateError || ""}
      </Snackbar>
      <Snackbar
        open={openUrlUpdateSuccess}
        severity="success"
        onClose={onUrlUpdateSuccessClose}
      >
        Successfully updated.
      </Snackbar>
    </>
  );
}

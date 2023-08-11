"use client";

import { TextField } from "@mui/material";
import Grid from "@mui/material/Unstable_Grid2";

import ProgressButton from "@/components/ProgressButton";
import Snackbar from "@/components/Snackbar";

import { useUrl } from "./hooks";

type Props = {
  name: string;
  url: string;
};

export default function UpdateForm({ name, url }: Props) {
  const {
    urlRef,
    urlUpdating,
    onUrlUpdate,
    openUrlUpdateSuccess,
    onUrlUpdateSuccessClose,
    urlUpdateError,
    onUrlUpdateErrorClose,
  } = useUrl(name);

  return (
    <>
      <Grid xs={12}>
        <TextField label="URL" inputRef={urlRef} fullWidth defaultValue={url} />
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

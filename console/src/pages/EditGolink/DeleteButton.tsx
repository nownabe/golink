import ProgressButton from "@/components/ProgressButton";
import { Golink } from "@/gen/golink/v1/golink_pb";
import Grid from "@mui/material/Unstable_Grid2";
import { useDeleteButton } from "./hooks";
import Snackbar from "@/components/Snackbar";

type Props = {
  golink: Golink;
};

export default function DeleteButton({ golink }: Props) {
  const {
    openSuccess,
    error,
    deleting,
    deleteGolink,
    onSuccessClose,
    onErrorClose,
    isOwner,
  } = useDeleteButton(golink);
  return (
    <>
      <Grid xs={12}>
        <ProgressButton
          loading={deleting}
          onClick={deleteGolink}
          disabled={!isOwner}
          color="error"
        >
          Delete
        </ProgressButton>
      </Grid>
      <Snackbar open={!!error} severity="error" onClose={onErrorClose}>
        {error || ""}
      </Snackbar>
      <Snackbar open={openSuccess} severity="success" onClose={onSuccessClose}>
        Successfully deleted.
      </Snackbar>
    </>
  );
}

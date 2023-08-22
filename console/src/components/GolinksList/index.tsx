import {
  Divider,
  List,
  ListItem,
  ListItemButton,
  ListItemText,
  Skeleton,
} from "@mui/material";

import LinkComponent from "@/components/LinkComponent";
import { Golink } from "@/gen/golink/v1/golink_pb";

type Props = {
  golinks: Golink[];
};

export function Loading() {
  return (
    <List>
      {Array.from({ length: 3 }).map((_, i) => (
        <div key={i}>
          <ListItem key={i} disablePadding>
            <ListItemButton sx={{ pl: 1 }}>
              <Skeleton
                variant="rectangular"
                sx={{ width: "100%", height: "56px" }}
              />
            </ListItemButton>
          </ListItem>
          <Divider />
        </div>
      ))}
    </List>
  );
}

export default function GolinksList({ golinks }: Props) {
  return (
    <List>
      {golinks.map((golink) => (
        <div key={golink.name}>
          <ListItem key={golink.name} disablePadding>
            <ListItemButton
              LinkComponent={LinkComponent}
              href={`/${golink.name}`}
              sx={{ pl: 1 }}
            >
              <ListItemText
                primary={`go/${golink.name}`}
                secondary={golink.url}
              />
            </ListItemButton>
          </ListItem>
          <Divider />
        </div>
      ))}
    </List>
  );
}

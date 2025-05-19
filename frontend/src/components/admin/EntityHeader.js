import styles from "./EntityHeader.module.css";
import { IconButton, Tooltip } from "@mui/material";
import { Add } from "@mui/icons-material";

const EntityHeader = ({ title, entityName, onAdd }) => {
  return (
    <div className={styles.top}>
      <h1>{title}</h1>
      <Tooltip title={`افزودن ${entityName}`} placement="right" arrow>
        <IconButton variant="addButton" onClick={onAdd}>
          <Add />
        </IconButton>
      </Tooltip>
    </div>
  );
};

export default EntityHeader;

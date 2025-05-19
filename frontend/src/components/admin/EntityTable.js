import styles from "./EntityTable.module.css";
import { IconButton } from "@mui/material";
import { Edit, Delete } from "@mui/icons-material";

const EntityTable = ({ data, columns, isTablet, onEdit, onDelete }) => {
  if (data.length === 0) {
    return <p>موردی برای نمایش وجود ندارد.</p>;
  }

  return (
    <table className={styles.table}>
      <tbody>
        {data.map((item) => (
          <tr key={item.id} className={styles.tr}>
            {columns.map(
              (col) =>
                (!isTablet || !col.hideOnTablet) && (
                  <td
                    key={col.key}
                    className={styles.td}
                    style={{ width: `${100 / (columns.length + 1)}%` }}
                  >
                    {col.render(item)}
                    {isTablet && col.customRender && col.customRender(item)}
                  </td>
                )
            )}
            <td className={styles.td}>
              <IconButton onClick={() => onEdit(item)}>
                <Edit />
              </IconButton>
            </td>
            <td className={styles.td}>
              <IconButton onClick={() => onDelete(item)}>
                <Delete />
              </IconButton>
            </td>
          </tr>
        ))}
      </tbody>
    </table>
  );
};

export default EntityTable;

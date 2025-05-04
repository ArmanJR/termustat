import styles from "./Universities.module.css";

import {
  IconButton,
  Tooltip,
  useTheme,
  useMediaQuery,
} from "@mui/material";

import {
  Add,
  Edit,
  Delete,
  Block,
} from "@mui/icons-material";

const universities = [
  { id: 1, name_en: "University 1", name_fa: "دانشگاه ۱", is_active: true },
  { id: 2, name_en: "University 2", name_fa: "دانشگاه ۲", is_active: true },
  { id: 3, name_en: "University 3", name_fa: "دانشگاه ۳", is_active: false },
];

const Universities = () => {
  const theme = useTheme();
  const isTablet = useMediaQuery(theme.breakpoints.down("md"));

  return (
    <div>
      <div className={styles.top}>
        <h1>دانشگاه‌ها</h1>
        <Tooltip title="افزدون دانشگاه" placement="right" arrow>
          <IconButton
            sx={{
              backgroundColor: "#309a9a",
              color: "#ffffff",
              "&:hover": {
                backgroundColor: "#42baba",
              },
            }}
          >
            <Add sx={{fontSize: "1.2em"}} />
          </IconButton>
        </Tooltip>
      </div>

      <table className={styles.table}>
        <tbody>
          {universities.map((uni) => (
            <tr key={uni.id} className={styles.tr}>
              <td className={styles.td}>
                { uni.name_fa }
                { (isTablet && !uni.is_active) && <IconButton disabled><Block /></IconButton> }
              </td>
              { !isTablet && <td className={styles.td}>{ uni.name_en }</td> }
              { !isTablet && <td className={styles.td}>{ uni.is_active ? "فعال" : "غیرفعال" }</td> }
              <td className={styles.td}><IconButton><Edit /></IconButton></td>
              <td className={styles.td}><IconButton><Delete /></IconButton></td>
            </tr>
          ))}
        </tbody>
      </table>
    </div>
  );
};

export default Universities;

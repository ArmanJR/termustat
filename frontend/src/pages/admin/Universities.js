import { useState, useEffect, useRef } from "react";
import styles from "./Universities.module.css";
import UniversitiesTheme from "./UniversitiesTheme";
import UniversityForm from "../../components/admin/UniversityForm";
import { getUniversities } from "../../api/admin/universities";

import {
  IconButton,
  Tooltip,
  useTheme,
  useMediaQuery,
  ThemeProvider,
} from "@mui/material";

import { Add, Edit, Delete, Block } from "@mui/icons-material";

const Universities = () => {
  const theme = useTheme();
  const isTablet = useMediaQuery(theme.breakpoints.down("md"));

  const [dialog, setDialog] = useState({
    open: false,
    mode: null,
    university: null,
  });

  const openDialog = (mode, university = null) => {
    setDialog({ open: true, mode, university });
  };

  const closeDialog = () => {
    setDialog((prev) => ({ ...prev, open: false }));
    setTimeout(() => {
      setDialog({ open: false, mode: null, university: null });
    }, 300);
  };

  const [universities, setUniversities] = useState([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState(null);
  const isFirstRender = useRef(true);

  const fetchUniversities = async () => {
    try {
      const response = await getUniversities();
      setUniversities(response.data);
    } catch (error) {
      if (error.response?.status === 500)
        setError("مشکلی در سرور رخ داده است. لطفا دوباره تلاش کنید.");
      else
        setError(error.message);
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    if (isFirstRender.current) {
      isFirstRender.current = false;
      return; // Skip the effect on first render (loading state)
    }
    fetchUniversities();
  }, []);

  return (
    <ThemeProvider theme={UniversitiesTheme}>
      <div>
        <div className={styles.top}>
          <h1>دانشگاه‌ها</h1>
          <Tooltip title="افزودن دانشگاه" placement="right" arrow>
            <IconButton variant="addButton" onClick={() => openDialog("add")}>
              <Add />
            </IconButton>
          </Tooltip>
        </div>

        {loading ? (
          <p></p>
        ) : error ? (
          <p>{error}</p>
        ) : universities.length === 0 ? (
          <p>دانشگاهی برای نمایش وجود ندارد.</p>
        ) : (
          <>
            <table className={styles.table}>
              <tbody>
                {universities.map((uni) => (
                  <tr key={uni.id} className={styles.tr}>
                    <td className={styles.td}>
                      {uni.name_fa}
                      {isTablet && !uni.is_active && (
                        <IconButton disabled>
                          <Block />
                        </IconButton>
                      )}
                    </td>
                    {!isTablet && <td className={styles.td}>{uni.name_en}</td>}
                    {!isTablet && (
                      <td className={styles.td}>
                        {uni.is_active ? "فعال" : "غیرفعال"}
                      </td>
                    )}
                    <td className={styles.td}>
                      <IconButton onClick={() => openDialog("edit", uni)}>
                        <Edit />
                      </IconButton>
                    </td>
                    <td className={styles.td}>
                      <IconButton onClick={() => openDialog("delete", uni)}>
                        <Delete />
                      </IconButton>
                    </td>
                  </tr>
                ))}
              </tbody>
            </table>
          </>
        )}
        <UniversityForm
          open={dialog.open}
          handleClose={closeDialog}
          university={dialog.university}
          mode={dialog.mode}
          refetchUniversities={fetchUniversities}
        />
      </div>
    </ThemeProvider>
  );
};

export default Universities;

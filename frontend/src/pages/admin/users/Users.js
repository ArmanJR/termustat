import { useState, useEffect } from "react";
import styles from "./Users.module.css";
import EntityPage from "../../../components/admin/EntityPage";
import config from "./config";
import useFetch from "../../../hooks/useFetch";
import { getUsers } from "../../../api/admin/users";
import { getUniversityById } from "../../../api/admin/universities";
import { getFacultyById } from "../../../api/admin/faculties";
import { useSearchParams } from "react-router-dom";
import { Pagination } from "@mui/material";

const enhanceUsersWithNames = async (users) => {
  const enhancedUsers = await Promise.all(
    users.map(async (user) => {
      let faculty_name = "نامشخص";
      let university_name = "نامشخص";
      if (user.faculty_id) {
        try {
          const res = await getFacultyById(user.faculty_id);
          faculty_name = res.data.name_fa;
        } catch {
          faculty_name = "نامشخص";
        }
      }
      if (user.university_id) {
        try {
          const res = await getUniversityById(user.university_id);
          university_name = res.data.name_fa;
        } catch {
          university_name = "نامشخص";
        }
      }
      return {
        ...user,
        faculty_name,
        university_name,
      };
    })
  );
  return enhancedUsers;
};

const Users = () => {
  const [searchParams, setSearchParams] = useSearchParams();
  const page = parseInt(searchParams.get("page")) || 1;

  const handlePageChange = (event, value) => {
    setSearchParams({ page: value });
  };

  const { data, loading, error, fetchData } = useFetch(() => getUsers(page), [page]);
  let users = data?.items;

  const [enhancedUsers, setEnhancedUsers] = useState([]);
  const [enhancing, setEnhancing] = useState(true);

  useEffect(() => {
    if (users) {
      if (users.length === 0) {
        setEnhancedUsers([]);
        setEnhancing(false);
      } else {
        setEnhancing(true);
        enhanceUsersWithNames(users).then((result) => {
          setEnhancedUsers(result);
          setEnhancing(false);
        });
      }
    }
  }, [users]);

  if (loading || enhancing) return <></>;
  if (error) return <div>{error}</div>;

  return (
    <>
      <EntityPage
        key={page}
        title={config.title}
        entityName={config.entityName}
        data={enhancedUsers}
        fetchData={fetchData}
        tableColumns={config.tableColumns}
        dialogFields={config.dialogFields}
        onSubmit={config.onSubmit}
        onOpenDialog={config.onOpenDialog}
        canAdd={false}
      />
      <div className={styles.pagination}>
        <Pagination
          count={Math.ceil(data.total / data.limit)}
          page={page}
          onChange={handlePageChange}
        />
      </div>
    </>
  );
};

export default Users;

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

const facultyCache = new Map();
const universityCache = new Map();

const enhanceUsersWithNames = async (users) => {
  const uniqueFacultyIds = [
    ...new Set(users.map((u) => u.faculty_id).filter(Boolean)),
  ];
  const uniqueUniversityIds = [
    ...new Set(users.map((u) => u.university_id).filter(Boolean)),
  ];
  const facultyPromises = uniqueFacultyIds
    .filter((id) => !facultyCache.has(id))
    .map(async (id) => {
      try {
        const res = await getFacultyById(id);
        facultyCache.set(id, res.data.name_fa);
      } catch(err) {
        console.error(`Failed to fetch faculty with ID ${id}:`, err);
        facultyCache.set(id, "نامشخص");
      }
    });
  const universityPromises = uniqueUniversityIds
    .filter((id) => !universityCache.has(id))
    .map(async (id) => {
      try {
        const res = await getUniversityById(id);
        universityCache.set(id, res.data.name_fa);
      } catch(err) {
        console.error(`Failed to fetch university with ID ${id}:`, err);
        universityCache.set(id, "نامشخص");
      }
    });
  await Promise.all([...facultyPromises, ...universityPromises]);
  return users.map((user) => ({
    ...user,
    faculty_name: facultyCache.get(user.faculty_id) || "نامشخص",
    university_name: universityCache.get(user.university_id) || "نامشخص",
  }));
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

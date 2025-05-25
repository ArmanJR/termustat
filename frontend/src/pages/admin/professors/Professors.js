import { useState, useEffect } from "react";
import styles from "./Professors.module.css";
import EntityTable from "../../../components/admin/EntityTable";
import useFetch from "../../../hooks/useFetch";
import { getUniversities } from "../../../api/admin/universities";
import { getProfessors } from "../../../api/admin/professors";
import config from "./config";

const Professors = () => {
  const { data: universities } = useFetch(getUniversities);
  const [selectedUniversity, setSelectedUniversity] = useState("");
  const { data: professors, error, fetchData } = useFetch(
    selectedUniversity ? () => getProfessors(selectedUniversity) : null
  );

  useEffect(() => {
    if (selectedUniversity) {
      fetchData();
    }
  }, [selectedUniversity]);

  return (
    <>
      <h1>اساتید</h1>

      <div className={styles.selectContainer}>
        <select
          className={styles.select}
          value={selectedUniversity}
          onChange={(e) => setSelectedUniversity(e.target.value)}
        >
          <option value="" disabled>
            دانشگاه را انتخاب کنید
          </option>
          {universities?.map((uni) => (
            <option key={uni.id} value={uni.id}>
              {uni.name_fa}
            </option>
          ))}
        </select>
      </div>

      {selectedUniversity &&
        (error ? (
          <p>{error}</p>
        ) : (
          professors && (
            <>
              <h1>
                {`اساتید ${universities?.find((u) => u.id === selectedUniversity)?.name_fa || ""}`}
              </h1>
              <EntityTable
                data={professors}
                columns={config.tableColumns}
              />
            </>
          )
        ))}
    </>
  );
};

export default Professors;

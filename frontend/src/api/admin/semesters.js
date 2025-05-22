import axiosInstance from "../axiosInstance";

export const getSemesters = async () => {
  const response = await axiosInstance.get(
    "/admin/semesters"
  );
  return response;
};

export const addSemester = async (semesterData) => {
  const response = await axiosInstance.post(
    "/admin/semesters",
    semesterData
  );
  return response;
};

export const deleteSemester = async (id) => {
  const response = await axiosInstance.delete(
    `/admin/semesters/${id}`
  );
  return response;
};

export const editSemester = async (id, semesterData) => {
  const response = await axiosInstance.put(
    `/admin/semesters/${id}`,
    semesterData
  );
  return response;
};

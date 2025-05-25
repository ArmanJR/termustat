import axiosInstance from "../axiosInstance";

export const getProfessors = async (universityId) => {
  const response = await axiosInstance.get(
    `/admin/universities/${universityId}/professors`
  );
  return response;
};

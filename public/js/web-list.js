let projectsData = [];
function addProject(e) {
  e.preventDefault();

  let projectName = document.getElementById("project-name").value;
  let startDate = document.getElementById("start-date").value;
  let endDate = document.getElementById("end-date").value;
  let description = document.getElementById("description").value;
  let technology1 = document.getElementById("technology-1").checked ? document.getElementById("technology-1").value : false;
  let technology2 = document.getElementById("technology-2").checked ? document.getElementById("technology-2").value : false;
  let technology3 = document.getElementById("technology-3").checked ? document.getElementById("technology-3").value : false;
  let technology4 = document.getElementById("technology-4").checked ? document.getElementById("technology-4").value : false;
  let uploadImage = document.getElementById("upload-image").files;

  if (projectName == "") 
  alert("You must fill in The Project Name field");
  else if (startDate == "") 
  alert("You must fill in The Start Date field");
  else if (endDate == "") 
  alert("You must fill in The End Date field");
  else if (description == "") 
  alert("You must fill in The Description field");
  else if (technology1 == false && technology2 == false && technology3 == false && technology4 == false) 
  alert("You must choose at least one Technologies field");
  else if 
  (uploadImage.length == 0) alert("You must upload The Project Image");
  else {
    uploadImage = URL.createObjectURL(uploadImage[0]);
    technology1 = technology1 != false ? `<i class="fa-brands ${technology1} fa-3x"></i>` : false;
    technology2 = technology2 != false ? `<i class="fa-brands ${technology2} fa-3x"></i>` : false;
    technology3 = technology3 != false ? `<i class="fa-brands ${technology3} fa-3x"></i>` : false;
    technology4 = technology4 != false ? `<i class="fa-brands ${technology4} fa-3x"></i>` : false;
    startDate = new Date(startDate);
    endDate = new Date(endDate);
    let duration = (endDate - startDate)/ (30*24*60*60*1000);
    duration = duration.toFixed(1);
    
    let projectData = {
      projectName,
      startDate,
      endDate,
      duration,
      description,
      technology1,
      technology2,
      technology3,
      technology4,
      uploadImage
    };
    projectsData.push(projectData);
    renderProject();
  }
}
function renderProject() {
  console.log(projectsData);
  document.getElementById("projects-container").innerHTML = "";
  for (let project of projectsData) {
    document.getElementById("projects-container").innerHTML += `
    <div id="project-container">
      <img src="${project.uploadImage}"/>
      <a href= " .../index.html " ><h3>${project.projectName}</h2></a>
      <div>durasi: ${project.duration} bulan</div>
      <p>
        ${project.description}
      </p>
      <div>
        ${project.technology1 != false ? project.technology1 : ""}
        ${project.technology2 != false ? project.technology2 : ""}
        ${project.technology3 != false ? project.technology3 : ""}
        ${project.technology4 != false ? project.technology4 : ""}
      </div>
      <div>
        <button>edit</button>
        <button>delete</button>
      </div>
    </div>
    `;
  }
}
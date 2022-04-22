'use strict';

// Request the backend api by functionName with args.
// @param timeout The timeout in ms, default to 300s.
async function srs_cloud_request(functionName, args, timeout) {
  try {
    return await new Promise((resolve, reject) => {
      $.ajax({
        type: 'POST',
        url: `/plugin?action=a&s=${functionName}&name=srs_cloud`,
        data: args,
        timeout: timeout || 300 * 1000,
        success: function (res, status, xhr) {
          if (!res?.status) {
            return reject({xhr, err: JSON.stringify(res)});
          }
          return resolve(JSON.parse(res.msg));
        },
        error: function (xhr, status, err) {
          reject({xhr, err});
        },
      });
    });
  } catch ({xhr, err}) {
    console.error(`Call ${functionName} with args=${JSON.stringify(args)} err`, status, err, xhr);
    layer.msg(`Call ${functionName} with args=${JSON.stringify(args)} err ${err}`, {icon: 2});
    throw xhr;
  }
}

async function install_dependence(deps) {
  const depsContent = deps.map((dep, index) => {
    return `
          <li>
            ${index + 1}: 
            <a href="javascript:void(0);" class="btlink" onclick="${dep.call}('${dep.id}')">
                安装${dep.title}
            </a>
          </li> 
    `;
  }).join('\n');

  layer.closeAll();
  return new Promise((resolve, reject) => {
    layer.open({
      icon: 0,
      closeBtn: 0,
      title: '安装云SRS',
      area: '400px',
      btn: [],
      content: `
        <p>云SRS依赖以下软件，请点击安装所需要的版本：</p>
        <ul>${depsContent}</ul>
        <p style="color: red; font-weight: bold;">请安装后所有依赖后，继续安装云SRS</p>
      `,
      success: function (layero, index) {
        // Note that this is called when the window is open, not closed or finished.
        resolve();
      },
    });
  });
}

async function installing_dependence(deps) {
  const depsContent = deps.map((dep, index) => {
    const desc = `正在安装<a href="javascript:void(0)" class="btlink" onclick="messagebox()">${dep.title}</a>`;
    return `<li>${index + 1}: ${desc}</li>`;
  }).join('\n');

  layer.closeAll();
  return new Promise((resolve, reject) => {
    layer.open({
      icon: 0,
      closeBtn: 0,
      title: '安装云SRS',
      area: '400px',
      btn: [],
      content: `
        <p>云SRS正在安装依赖软件：</p>
        <ul>${depsContent}</ul>
        请耐心等待安装完成，请点<a href="javascript:void(0)" class="btlink" onclick="messagebox()">这里</a>查看安装详情
      `,
      success: function (layero, index) {
        // Note that this is called when the window is open, not closed or finished.
        resolve();
      },
    });
  });
}


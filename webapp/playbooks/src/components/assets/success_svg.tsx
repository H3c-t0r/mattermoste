// Copyright (c) 2015-present Mattermost, Inc. All Rights Reserved.
// See LICENSE.txt for license information.

import React from 'react';
import styled from 'styled-components';

import Icon from 'src/components/assets/svg';

const Svg = styled(Icon)`
    width: 256px;
    height: 156px;
    margin-right: 54px;
    margin-left: 54px;
`;

const SuccessSvg = () => (
    <Svg
        viewBox='0 0 256 156'
        fill='none'
        xmlns='http://www.w3.org/2000/svg'
    >
        <path
            opacity='0.1'
            d='M154.728 52.6595C153.343 52.3978 151.911 52.5535 150.616 53.1065C149.834 53.4517 148.988 53.6282 148.133 53.6244C147.278 53.6206 146.434 53.4366 145.655 53.0846C144.737 52.6837 143.742 52.4861 142.739 52.5055C141.737 52.5249 140.751 52.7608 139.848 53.1969C139.333 53.4694 138.759 53.6128 138.176 53.6146C135.822 53.6146 133.863 51.2525 133.456 48.1357C133.926 47.7927 134.325 47.3632 134.633 46.8703C136.013 44.6548 138.152 43.2332 140.552 43.2332C142.953 43.2332 145.064 44.6377 146.437 46.8313C146.85 47.4947 147.429 48.0404 148.116 48.4154C148.804 48.7903 149.577 48.9818 150.361 48.9711H150.422C152.303 48.9808 153.926 50.4782 154.728 52.6595Z'
            fill='var(--button-bg)'
        />
        <path
            opacity='0.1'
            d='M218.4 76.0438C216.577 75.7026 214.693 75.9056 212.987 76.6267C211.958 77.0767 210.845 77.3068 209.719 77.3019C208.594 77.297 207.482 77.057 206.458 76.5979C205.249 76.0753 203.939 75.8178 202.62 75.8431C201.3 75.8684 200.002 76.1759 198.815 76.7444C198.137 77.0997 197.381 77.2867 196.613 77.2891C193.515 77.2891 190.936 74.2094 190.4 70.1458C191.019 69.6987 191.545 69.1388 191.949 68.4962C193.766 65.6077 196.581 63.7542 199.741 63.7542C202.901 63.7542 205.68 65.5853 207.487 68.4451C208.031 69.3101 208.792 70.0216 209.697 70.5105C210.602 70.9993 211.62 71.249 212.652 71.235H212.732C215.208 71.2478 217.345 73.1999 218.4 76.0438Z'
            fill='var(--button-bg)'
        />
        <path
            opacity='0.1'
            d='M161.888 42.8765L158.082 45.2826L160.392 41.0958C159.737 40.5763 158.927 40.2884 158.09 40.2775H158.029C157.764 40.2822 157.5 40.2634 157.239 40.2213L155.954 41.0348L156.508 40.0333C155.594 39.7153 154.804 39.1166 154.252 38.3234L151.942 39.789L153.414 37.1484C152.063 35.5338 150.243 34.5421 148.242 34.5421C145.842 34.5421 143.703 35.9662 142.323 38.1817C141.915 38.8469 141.336 39.3921 140.647 39.7612C139.957 40.1304 139.182 40.3102 138.399 40.2824H138.269C135.619 40.2824 133.471 43.2771 133.471 46.968C133.471 50.659 135.619 53.6537 138.269 53.6537C138.852 53.6518 139.426 53.5093 139.942 53.2384C140.844 52.8013 141.83 52.5647 142.832 52.5448C143.835 52.525 144.83 52.7225 145.748 53.1236C146.525 53.471 147.365 53.6516 148.216 53.6537C149.067 53.6558 149.909 53.4794 150.687 53.1358C151.599 52.7424 152.584 52.5486 153.578 52.5675C154.571 52.5865 155.549 52.8178 156.444 53.2458C156.955 53.5117 157.521 53.6516 158.097 53.6537C160.748 53.6537 162.896 50.6614 162.896 46.968C162.91 45.5415 162.564 44.1344 161.888 42.8765Z'
            fill='var(--button-bg)'
        />
        <path
            opacity='0.1'
            d='M228.257 62.9214L223.187 66.132L226.264 60.5453C225.391 59.8521 224.312 59.468 223.197 59.4534H223.115C222.763 59.4596 222.411 59.4346 222.063 59.3785L220.352 60.4639L221.09 59.1275C219.872 58.7032 218.819 57.9044 218.085 56.8459L215.007 58.8015L216.967 55.2781C215.168 53.1237 212.744 51.8003 210.078 51.8003C206.88 51.8003 204.032 53.7006 202.193 56.6568C201.649 57.5445 200.878 58.2719 199.96 58.7645C199.041 59.2571 198.008 59.4971 196.966 59.4599H196.793C193.262 59.4599 190.4 63.456 190.4 68.381C190.4 73.3059 193.262 77.302 196.793 77.302C197.569 77.2995 198.334 77.1094 199.021 76.7479C200.222 76.1646 201.536 75.8488 202.871 75.8224C204.207 75.7959 205.533 76.0594 206.756 76.5947C207.79 77.0582 208.91 77.2992 210.044 77.302C211.177 77.3048 212.299 77.0693 213.335 76.611C214.549 76.086 215.863 75.8274 217.186 75.8527C218.509 75.878 219.812 76.1866 221.005 76.7577C221.685 77.1125 222.44 77.2992 223.207 77.302C226.738 77.302 229.599 73.3092 229.599 68.381C229.619 66.4775 229.157 64.5998 228.257 62.9214Z'
            fill='var(--button-bg)'
        />
        <path
            d='M127.166 35.4361C127.137 35.2132 127.068 34.9974 126.963 34.7986C126.879 34.68 126.779 34.5746 126.664 34.4859C126.479 34.3143 126.262 34.1799 126.026 34.0902C125.906 34.0435 125.776 34.0267 125.648 34.0411C125.52 34.0556 125.397 34.1009 125.29 34.1732C125.101 34.3433 124.986 34.5813 124.972 34.8352C124.92 35.1795 124.98 35.5313 125.143 35.8391C125.34 36.1787 125.678 36.4059 125.96 36.677C126.167 36.8728 126.347 37.0955 126.494 37.339C126.562 37.4661 126.644 37.585 126.74 37.6932C126.787 37.7467 126.847 37.7883 126.914 37.8146C126.981 37.8409 127.053 37.8512 127.125 37.8446C127.244 37.8138 127.353 37.7504 127.439 37.6614C127.554 37.5612 127.868 37.378 127.877 37.2168C127.887 37.0556 127.632 36.7918 127.581 36.6697C127.399 36.2741 127.261 35.8607 127.166 35.4361Z'
            fill='var(--online-indicator)'
        />
        <path
            d='M128 156C198.692 156 256 153.485 256 150.382C256 147.279 198.692 144.764 128 144.764C57.3076 144.764 0 147.279 0 150.382C0 153.485 57.3076 156 128 156Z'
            fill='var(--center-channel-color)'
            fillOpacity='0.08'
        />
        <path
            d='M216.08 148.525H205.663C205.09 148.525 204.541 148.299 204.135 147.895C203.73 147.491 203.503 146.944 203.503 146.373C203.503 145.803 203.73 145.255 204.135 144.852C204.541 144.448 205.09 144.221 205.663 144.221C205.373 144.231 205.085 144.183 204.814 144.079C204.544 143.975 204.297 143.818 204.089 143.618C203.88 143.417 203.715 143.177 203.601 142.911C203.488 142.645 203.43 142.359 203.43 142.071C203.43 141.782 203.488 141.496 203.601 141.23C203.715 140.964 203.88 140.724 204.089 140.523C204.297 140.323 204.544 140.166 204.814 140.062C205.085 139.958 205.373 139.91 205.663 139.92H206.139C206.411 138.068 208.507 136.627 211.058 136.625C213.608 136.622 215.717 138.061 215.984 139.91H216.08C216.369 139.9 216.658 139.949 216.928 140.052C217.199 140.156 217.446 140.313 217.654 140.514C217.862 140.714 218.028 140.955 218.141 141.221C218.254 141.486 218.313 141.772 218.313 142.061C218.313 142.349 218.254 142.635 218.141 142.901C218.028 143.167 217.862 143.407 217.654 143.608C217.446 143.809 217.199 143.965 216.928 144.069C216.658 144.173 216.369 144.221 216.08 144.212C216.653 144.212 217.202 144.438 217.607 144.842C218.012 145.245 218.24 145.793 218.24 146.364C218.24 146.934 218.012 147.482 217.607 147.885C217.202 148.289 216.653 148.516 216.08 148.516V148.525Z'
            fill='var(--online-indicator)'
        />
        <path
            opacity='0.1'
            d='M218.238 142.069C218.237 142.639 218.009 143.186 217.605 143.589C217.2 143.993 216.652 144.22 216.08 144.221H205.666L216.013 139.91H216.077C216.651 139.911 217.201 140.139 217.606 140.544C218.011 140.949 218.238 141.498 218.238 142.069Z'
            fill='var(--center-channel-color)'
            fillOpacity='0.48'
        />
        <path
            opacity='0.1'
            d='M203.567 146.395C203.569 146.965 203.797 147.511 204.202 147.914C204.606 148.317 205.155 148.544 205.727 148.545H216.143L205.788 144.253H205.732C205.16 144.254 204.612 144.479 204.207 144.88C203.801 145.282 203.571 145.826 203.567 146.395Z'
            fill='var(--center-channel-color)'
            fillOpacity='0.48'
        />
        <path
            d='M152.766 145.956H146.854C146.529 145.956 146.217 145.827 145.987 145.598C145.757 145.369 145.628 145.058 145.628 144.734C145.628 144.41 145.757 144.1 145.987 143.871C146.217 143.642 146.529 143.513 146.854 143.513C146.529 143.513 146.217 143.384 145.987 143.155C145.757 142.926 145.628 142.616 145.628 142.292C145.628 141.968 145.757 141.657 145.987 141.428C146.217 141.199 146.529 141.07 146.854 141.07H147.126C147.278 140.02 148.47 139.202 149.917 139.202C151.364 139.202 152.555 140.015 152.712 141.065H152.766C153.091 141.065 153.403 141.194 153.633 141.423C153.863 141.652 153.992 141.963 153.992 142.287C153.992 142.611 153.863 142.921 153.633 143.15C153.403 143.379 153.091 143.508 152.766 143.508C153.091 143.508 153.403 143.637 153.633 143.866C153.863 144.095 153.992 144.405 153.992 144.729C153.992 145.053 153.863 145.364 153.633 145.593C153.403 145.822 153.091 145.951 152.766 145.951V145.956Z'
            fill='var(--online-indicator)'
        />
        <path
            opacity='0.1'
            d='M153.992 142.292C153.992 142.616 153.863 142.926 153.633 143.155C153.403 143.384 153.091 143.513 152.766 143.513H146.854L152.739 141.07H152.766C153.091 141.07 153.403 141.199 153.633 141.428C153.863 141.657 153.992 141.968 153.992 142.292Z'
            fill='var(--center-channel-color)'
            fillOpacity='0.48'
        />
        <path
            opacity='0.1'
            d='M145.667 144.742C145.667 145.066 145.797 145.376 146.027 145.605C146.257 145.834 146.568 145.963 146.894 145.963H152.803L146.918 143.52H146.881C146.558 143.523 146.25 143.654 146.022 143.882C145.795 144.111 145.667 144.42 145.667 144.742Z'
            fill='var(--center-channel-color)'
            fillOpacity='0.48'
        />
        <path
            d='M42.0414 152.886H29.5356C29.1887 152.896 28.8432 152.837 28.5195 152.712C28.1959 152.588 27.9007 152.399 27.6516 152.159C27.4024 151.918 27.2043 151.63 27.069 151.311C26.9337 150.993 26.864 150.651 26.864 150.305C26.864 149.959 26.9337 149.617 27.069 149.299C27.2043 148.98 27.4024 148.692 27.6516 148.451C27.9007 148.211 28.1959 148.022 28.5195 147.897C28.8432 147.772 29.1887 147.714 29.5356 147.724C29.1887 147.735 28.8432 147.676 28.5195 147.551C28.1959 147.426 27.9007 147.238 27.6516 146.997C27.4024 146.756 27.2043 146.468 27.069 146.15C26.9337 145.831 26.864 145.489 26.864 145.144C26.864 144.798 26.9337 144.456 27.069 144.137C27.2043 143.819 27.4024 143.531 27.6516 143.29C27.9007 143.049 28.1959 142.861 28.5195 142.736C28.8432 142.611 29.1887 142.552 29.5356 142.563H30.1045C30.4306 140.342 32.9489 138.613 36.0067 138.611C39.0645 138.608 41.5853 140.333 41.9163 142.551H42.0291C42.7024 142.571 43.3412 142.852 43.8101 143.334C44.279 143.816 44.5413 144.46 44.5413 145.131C44.5413 145.802 44.279 146.447 43.8101 146.929C43.3412 147.41 42.7024 147.691 42.0291 147.712C42.7024 147.733 43.3412 148.014 43.8101 148.495C44.279 148.977 44.5413 149.622 44.5413 150.293C44.5413 150.964 44.279 151.608 43.8101 152.09C43.3412 152.572 42.7024 152.853 42.0291 152.873L42.0414 152.886Z'
            fill='var(--online-indicator)'
        />
        <path
            opacity='0.1'
            d='M44.6284 145.142C44.6271 145.826 44.3543 146.481 43.8695 146.965C43.3848 147.449 42.7276 147.722 42.0414 147.724H29.5356L41.9482 142.551H42.0242C42.3661 142.55 42.7047 142.617 43.0206 142.747C43.3365 142.877 43.6236 143.068 43.8655 143.309C44.1073 143.549 44.2992 143.835 44.4301 144.15C44.561 144.464 44.6284 144.802 44.6284 145.142Z'
            fill='var(--center-channel-color)'
            fillOpacity='0.48'
        />
        <path
            opacity='0.1'
            d='M27.0271 150.321C27.0303 151.005 27.3049 151.659 27.791 152.142C28.277 152.625 28.935 152.897 29.6214 152.898H42.1272L29.7023 147.749H29.6263C28.9401 147.749 28.282 148.02 27.7951 148.502C27.3083 148.983 27.0323 149.637 27.0271 150.321Z'
            fill='var(--center-channel-color)'
            fillOpacity='0.48'
        />
        <path
            opacity='0.1'
            d='M67.6781 37.3878C66.7052 37.2018 65.6986 37.309 64.7871 37.6956C64.2432 37.9379 63.6541 38.0631 63.0583 38.0631C62.4626 38.0631 61.8735 37.9379 61.3296 37.6956C60.6866 37.4151 59.9901 37.2768 59.2882 37.2903C58.5863 37.3038 57.8957 37.4687 57.264 37.7738C56.9037 37.9642 56.5023 38.0648 56.0944 38.0669C54.4466 38.0669 53.0758 36.4132 52.7914 34.2294C53.1205 33.9901 53.4002 33.6899 53.6153 33.3452C54.5814 31.7941 56.0674 30.7999 57.7569 30.7999C59.4464 30.7999 60.9177 31.777 61.8838 33.3183C62.1727 33.7807 62.5762 34.1614 63.0554 34.4236C63.5347 34.6858 64.0736 34.8208 64.6203 34.8157H64.6645C65.9764 34.8132 67.1117 35.8612 67.6781 37.3878Z'
            fill='var(--button-bg)'
        />
        <path
            opacity='0.1'
            d='M72.6854 30.541L70.0224 32.224L71.6384 29.2928C71.1809 28.9273 70.6137 28.7252 70.0273 28.7187H69.9857C69.8009 28.723 69.6162 28.7099 69.4339 28.6796L68.5291 29.2488L68.9165 28.5477C68.2771 28.325 67.7245 27.907 67.3374 27.3533L65.7165 28.3767L66.7366 26.5276C65.7926 25.3991 64.5175 24.7054 63.1173 24.7054C61.4376 24.7054 59.9296 25.6996 58.9733 27.2507C58.6866 27.7135 58.2823 28.0927 57.8014 28.35C57.3205 28.6073 56.7799 28.7337 56.2342 28.7163H56.1533C54.2995 28.7163 52.7939 30.8121 52.7939 33.3965C52.7939 35.9808 54.2995 38.0767 56.1533 38.0767C56.5614 38.0758 56.963 37.9751 57.323 37.7835C57.9547 37.4785 58.6453 37.3135 59.3472 37.3C60.0491 37.2866 60.7455 37.4248 61.3886 37.7054C61.9357 37.9487 62.5279 38.0752 63.1271 38.0767C63.7185 38.0759 64.3032 37.952 64.8436 37.7127C65.4816 37.4375 66.1714 37.302 66.8665 37.3155C67.5616 37.329 68.2457 37.4911 68.8724 37.7909C69.2296 37.9773 69.6266 38.0754 70.0298 38.0767C71.886 38.0767 73.3892 35.9808 73.3892 33.3965C73.3982 32.401 73.1562 31.4191 72.6854 30.541Z'
            fill='var(--button-bg)'
        />
        <path
            opacity='0.1'
            d='M32.5198 64.2575C33.4888 64.074 34.4907 64.1811 35.3986 64.5653C35.942 64.8077 36.5307 64.933 37.1261 64.933C37.7215 64.933 38.3102 64.8077 38.8536 64.5653C39.4966 64.2848 40.1931 64.1465 40.895 64.16C41.5969 64.1735 42.2875 64.3384 42.9192 64.6435C43.2796 64.8337 43.6809 64.9342 44.0888 64.9366C45.7366 64.9366 47.1098 63.2829 47.3943 61.0991C47.0646 60.8605 46.7847 60.5602 46.5704 60.2149C45.6042 58.6638 44.106 57.6696 42.4263 57.6696C40.7466 57.6696 39.268 58.6467 38.3018 60.188C38.0118 60.6514 37.6066 61.0325 37.1255 61.2943C36.6444 61.5562 36.1037 61.69 35.5555 61.6829H35.5138C34.2093 61.6829 33.0715 62.7308 32.5198 64.2575Z'
            fill='var(--button-bg)'
        />
        <path
            opacity='0.1'
            d='M27.498 57.4106L30.1609 59.0936L28.5425 56.1624C29.0002 55.7972 29.5673 55.5952 30.1536 55.5883H30.1977C30.3825 55.5926 30.5672 55.5795 30.7494 55.5492L31.6518 56.1184L31.2644 55.4173C31.9045 55.1947 32.458 54.7768 32.846 54.2229L34.4619 55.2439L33.4394 53.3948C34.3859 52.2663 35.6586 51.5725 37.0612 51.5725C38.7409 51.5725 40.2366 52.5667 41.2052 54.1178C41.4917 54.5808 41.8959 54.9601 42.3769 55.2175C42.8579 55.4749 43.3986 55.6011 43.9442 55.5834H44.0325C45.8888 55.5834 47.3919 57.6793 47.3919 60.2636C47.3919 62.848 45.8888 64.9438 44.0325 64.9438C43.6245 64.9423 43.223 64.8417 42.8629 64.6507C42.2312 64.3457 41.5406 64.1807 40.8387 64.1672C40.1368 64.1537 39.4403 64.292 38.7973 64.5725C38.251 64.8159 37.6596 64.9424 37.0612 64.9438C36.4697 64.9436 35.885 64.8196 35.3447 64.5799C34.7063 64.3046 34.016 64.1691 33.3205 64.1826C32.625 64.1961 31.9406 64.3582 31.3134 64.658C30.9564 64.8448 30.5593 64.9428 30.156 64.9438C28.3022 64.9438 26.7966 62.848 26.7966 60.2636C26.7873 59.2692 27.0284 58.2882 27.498 57.4106Z'
            fill='var(--button-bg)'
        />
        <path
            d='M128 152C150.091 152 168 151.105 168 150C168 148.895 150.091 148 128 148C105.909 148 88 148.895 88 150C88 151.105 105.909 152 128 152Z'
            fill='black'
            fillOpacity='0.2'
        />
        <path
            d='M125 125.544C152.998 125.544 175.696 103.104 175.696 75.4239C175.696 47.7436 152.998 25.3043 125 25.3043C97.0016 25.3043 74.3043 47.7436 74.3043 75.4239C74.3043 103.104 97.0016 125.544 125 125.544Z'
            fill='white'
        />
        <path
            d='M178 76.0017C178 86.4841 174.891 96.731 169.067 105.447C163.243 114.162 154.965 120.955 145.281 124.966C135.596 128.977 124.939 130.027 114.659 127.981C104.378 125.936 94.9341 120.888 87.5221 113.475C80.1102 106.063 75.0628 96.6192 73.0181 86.3381C70.9734 76.057 72.0233 65.4005 76.0351 55.7162C80.0468 46.0318 86.8402 37.7546 95.5562 31.9312C104.272 26.1077 114.519 22.9997 125.002 23C131.962 23.0002 138.854 24.3713 145.284 27.035C151.714 29.6988 157.557 33.6029 162.478 38.5246C167.399 43.4462 171.303 49.289 173.966 55.7194C176.63 62.1497 178 69.0416 178 76.0017ZM118.875 104.056L158.19 64.7408C158.508 64.4234 158.761 64.0465 158.933 63.6314C159.105 63.2164 159.194 62.7714 159.194 62.3221C159.194 61.8727 159.105 61.4278 158.933 61.0127C158.761 60.5977 158.508 60.2207 158.19 59.9034L153.353 55.1014C153.036 54.7834 152.66 54.5311 152.245 54.359C151.831 54.1868 151.386 54.0982 150.937 54.0982C150.489 54.0982 150.044 54.1868 149.63 54.359C149.215 54.5311 148.839 54.7834 148.522 55.1014L116.457 87.1675L101.485 72.1601C100.842 71.5209 99.9724 71.1622 99.066 71.1622C98.1596 71.1622 97.29 71.5209 96.6473 72.1601L91.81 77.0259C91.1708 77.6687 90.8121 78.5382 90.8121 79.4446C90.8121 80.351 91.1708 81.2206 91.81 81.8633L114.039 104.084C114.356 104.402 114.733 104.655 115.148 104.827C115.563 104.999 116.008 105.088 116.458 105.088C116.907 105.088 117.352 104.999 117.767 104.827C118.182 104.655 118.559 104.402 118.876 104.084L118.875 104.056Z'
            fill='var(--online-indicator)'
        />
    </Svg>
);

export default SuccessSvg;

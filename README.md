# DeBox ChatWidget Demo 

Welcome to the **DeBox ChatWidget** demo repository. This repository demonstrates the usage of the ChatWidget in different environments and provides simple usage instructions.

### Usage Guide
For detailed usage instructions and documentation, please visit: [ChatWidget Usage Guide](https://docs.debox.pro/zh/ChatWidget).

### File Structure
The file structure of this repository is as follows:

```bash
/project-root
├── /html-umd
├── /html-npm
├── /html-react
└── /conversationId-server
```

- **html-umd**: UMD module example for usage in the browser.
- **html-npm**: Example using npm for installation and integration.
- **html-react**: Example for using the widget in a React environment.
- **conversationId-server**: Backend example demonstrating how to integrate with the ChatWidget.

### Different Usage Methods

1. **HTML (UMD)**
   This example shows how to use the ChatWidget with the UMD (Universal Module Definition) format in a native HTML environment. 
   Access the demo here:  
   [HTML UMD Demo](https://debox-chat-widget-demo.vercel.app/html-umd/index.html?project_id=12345&group_name=box&chain_id=1&contract_address=0x32b77729cd87f1ef2bea4c650c16f89f08472c69).

   The official umd file can be downloaded here:
   \[[index.umd.js](https://github.com/debox-pro/debox-ChatWidget-demo/blob/main/html-umd/umd-file/index.umd.js)\] \[[License](https://github.com/debox-pro/debox-ChatWidget-demo/blob/main/html-umd/umd-file/index.umd.js.LICENSE.txt)\]

2. **HTML (npm)**
   This example demonstrates how to integrate the ChatWidget using npm package manager in a native HTML environment.  
   Access the demo here:  
   [HTML NPM Demo](https://debox-chat-widget-demo.vercel.app/html-npm/index.html?project_id=12345&group_name=box&chain_id=1&contract_address=0x32b77729cd87f1ef2bea4c650c16f89f08472c69).

3. **React**
   This example shows how to use the ChatWidget in a React environment. Make sure to install the necessary dependencies via npm or yarn and configure according to the documentation.  
   Access the demo here:  
   [REACT NPM Demo](https://debox-chat-widget-demo.vercel.app/react-npm/index.html?project_id=12345&group_name=box&chain_id=1&contract_address=0x32b77729cd87f1ef2bea4c650c16f89f08472c69).

   
### conversationId-server Example

In the `conversationId-server` folder, you will find a simple backend example demonstrating how to interact with the DeBox Open Platform's `/openapi/chatwidget/conversation/id` API.

#### Why Use a Backend?

To use the ChatWidget, you need to retrieve the chat group ID by calling the DeBox Open Platform API. This API requires sensitive credentials such as `X-API-KEY` and `app_secret`. For security reasons, it is recommended to handle this API call on a private backend server.

In this example, a private backend API endpoint is created to securely call the DeBox Open Platform API using the private `app_secret`. The frontend can then use this private endpoint to retrieve the necessary data without exposing sensitive credentials.


---

We hope you have a great experience using the ChatWidget. If you encounter any issues or have any suggestions, please feel free to contact us.

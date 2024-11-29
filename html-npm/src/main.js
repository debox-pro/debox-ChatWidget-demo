import { DeBoxChatWidget } from '@debox-pro/chat-widget-html';

function getQueryVariable(variable) {
    return new URL(window.location.href).searchParams.get(variable);
}

// Initialize plugin
DeBoxChatWidget.init({
    projectId: getQueryVariable('project_id') ?? import.meta.env.VITE_PKEY,
    zIndex: '9999',
    containerDomId: 'test-dom',
});

// Check for input parameters in the URL and initialize Conversation ID if present
async function initializeConversationId() {
    const groupName = getQueryVariable('group_name');
    const chainId = getQueryVariable('chain_id');
    const contractAddress = getQueryVariable('contract_address');

    if (groupName || chainId || contractAddress) {
        const groupParam = groupName ? encodeURIComponent(groupName) : '';
        const chainParam = chainId ? encodeURIComponent(chainId) : '';
        const addressParam = contractAddress ? encodeURIComponent(contractAddress) : '';

        const url = `https://demo-chat.debox.pro/getconversationid?group_name=${groupParam}&chain_id=${chainParam}&contract_address=${addressParam}`;

        try {
            const response = await fetch(url);
            const data = await response.json();
            if (data.success) {
                const conversationId = data.data.cid;
                DeBoxChatWidget.setConversation(conversationId);
            } else {
                console.warn("Failed to retrieve Conversation ID during initialization");
            }
        } catch (error) {
            console.error("Error fetching Conversation ID during initialization", error);
        }
    }
}

// Execute initialization logic
initializeConversationId();

// Get Conversation ID based on inputs
async function getConversationId() {
    const groupName = document.getElementById('groupName').value;
    const chainId = document.getElementById('chainId').value;
    const contractAddress = document.getElementById('contractAddress').value;

    const groupParam = groupName ? encodeURIComponent(groupName) : '';
    const chainParam = chainId ? encodeURIComponent(chainId) : '';
    const addressParam = contractAddress ? encodeURIComponent(contractAddress) : '';

    if (!groupParam && !chainParam && !addressParam) {
        alert("Please fill in at least one field.");
        return;
    }

    const url = `https://demo-chat.debox.pro/getconversationid?group_name=${groupParam}&chain_id=${chainParam}&contract_address=${addressParam}`;

    try {
        const response = await fetch(url);
        const data = await response.json();
        if (data.success) {
            const conversationId = data.data.cid;
            document.getElementById('conversationInput').value = conversationId;
            alert(`Conversation ID: ${conversationId}`);
        } else {
            alert("Failed to get Conversation ID");
        }
    } catch (error) {
        alert("Error fetching Conversation ID");
        console.error(error);
    }
}

// Change conversation based on input
function setConversationFromInput() {
    const inputValue = document.getElementById('conversationInput').value;
    if (inputValue) {
        DeBoxChatWidget.setConversation(inputValue);
    } else {
        alert('Please enter a valid Conversation ID');
    }
}



// Listen for plugin events
    DeBoxChatWidget.addEventListener((e) => {
        const events = document.getElementById('events');
        const li = document.createElement('li');
        li.innerHTML = JSON.stringify(e.detail);
        events.appendChild(li);
    });


// Expose functions to global window object
window.getConversationId = getConversationId;
window.setConversationFromInput = setConversationFromInput;
window.DeBoxChatWidget = DeBoxChatWidget;




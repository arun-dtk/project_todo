document.addEventListener('DOMContentLoaded', function () {
    const signupForm = document.getElementById('signupForm');
    const loginForm = document.getElementById('loginForm');
    const todosList = document.getElementById('todosList');
    const addTodoBtn = document.getElementById('addTodoBtn');

    // Handle Signup Form Submission
    if (signupForm) {
        signupForm.addEventListener('submit', async function (e) {
            e.preventDefault();
            const email = document.getElementById('email').value;
            const password = document.getElementById('password').value;
            const firstname = document.getElementById('firstname').value;
            const lastname = document.getElementById('lastname').value;

            try {
                const response = await fetch('/signup', {
                    method: 'POST',
                    headers: {
                        'Content-Type': 'application/json',
                    },
                    body: JSON.stringify({ email, password, firstName: firstname, lastName: lastname }),
                });

                const data = await response.json();
                document.getElementById('message').textContent = data.message;
            } catch (error) {
                console.error('Error:', error);
            }
        });
    }

    // Handle Login Form Submission
    if (loginForm) {
        loginForm.addEventListener('submit', async function (e) {
            e.preventDefault();
            const email = document.getElementById('loginEmail').value;
            const password = document.getElementById('loginPassword').value;

            try {
                const response = await fetch('/login', {
                    method: 'POST',
                    headers: {
                        'Content-Type': 'application/json',
                    },
                    body: JSON.stringify({ email, password }),
                });

                const data = await response.json();
                document.getElementById('message').textContent = data.message;

                if (response.ok) {
                    localStorage.setItem('token', data.token);
                    window.location.href = '/todolist';
                }
            } catch (error) {
                console.error('Error:', error);
            }
        });
    }

    // Handle Todos Page
    if (todosList) {
        addTodoBtn.addEventListener('click', function () {
            renderTodoForm();
        });

        loadTodos();
    }

    async function loadTodos() {
        const token = localStorage.getItem('token');
        const response = await fetch('/todos', {
            headers: {
                'Authorization': `${token}`,
            },
        });
        const todos = await response.json();
        todos.forEach(todo => renderTodoItem(todo));
    }

    function renderTodoItem(todo) {
        const li = document.createElement('li');
        li.className = 'todo-item';

        const todoHeader = document.createElement('div');
        todoHeader.className = 'todo-header';

        const title = document.createElement('span');
        title.textContent = todo.title;

        const actions = document.createElement('div');
        actions.className = 'todo-actions';

        const editBtn = document.createElement('button');
        editBtn.textContent = 'Edit';
        editBtn.addEventListener('click', function () {
            renderTodoForm(todo);
        });

        const deleteBtn = document.createElement('button');
        deleteBtn.textContent = 'Delete';
        deleteBtn.addEventListener('click', function () {
            deleteTodo(todo.id);
        });

        actions.appendChild(editBtn);
        actions.appendChild(deleteBtn);
        todoHeader.appendChild(title);
        todoHeader.appendChild(actions);
        li.appendChild(todoHeader);

        todo.list.forEach(item => {
            const checklistItem = document.createElement('div');
            checklistItem.className = 'checklist-item';

            const checkbox = document.createElement('input');
            checkbox.type = 'checkbox';
            checkbox.checked = item.checked;
            checkbox.disabled = true;

            const itemText = document.createElement('span');
            itemText.textContent = item.item;

            checklistItem.appendChild(checkbox);
            checklistItem.appendChild(itemText);
            li.appendChild(checklistItem);
        });

        todosList.appendChild(li);
    }

    function renderTodoForm(todo = { title: '', list: [] }) {
        todosList.innerHTML = '';

        const formTemplate = document.getElementById('todoFormTemplate');
        const formClone = formTemplate.cloneNode(true);
        formClone.style.display = 'block';

        const titleInput = formClone.querySelector('#todoTitle');
        titleInput.value = todo.title;

        const checklistItems = formClone.querySelector('#checklistItems');
        todo.list.forEach(item => {
            addChecklistItem(checklistItems, item.item, item.checked);
        });

        const newChecklistItem = formClone.querySelector('#newChecklistItem');
        newChecklistItem.addEventListener('keyup', function (e) {
            if (e.key === 'Enter' && newChecklistItem.value.trim()) {
                addChecklistItem(checklistItems, newChecklistItem.value.trim());
                newChecklistItem.value = ''; // Clear the textarea
            }
        });

        const saveTodoBtn = formClone.querySelector('#saveTodoBtn');
        saveTodoBtn.addEventListener('click', function () {
            saveTodo({
                id: todo.id,
                title: titleInput.value,
                list: getChecklistItems(checklistItems),
            });
        });

        const backBtn = formClone.querySelector('#backBtn');  // Back button event listener
        backBtn.addEventListener('click', function () {
            todosList.innerHTML = '';
            loadTodos();  // Return to the main todo list view
        });
    

        todosList.appendChild(formClone);
    }

    function addChecklistItem(container, item = '', checked = false) {
        const li = document.createElement('li');
        li.className = 'checklist-item';

        const checkbox = document.createElement('input');
        checkbox.type = 'checkbox';
        checkbox.checked = checked;
        li.appendChild(checkbox);

        const input = document.createElement('input');
        input.type = 'text';
        input.value = item;
        li.appendChild(input);

        container.appendChild(li);
    }

    function getChecklistItems(container) {
        const items = [];
        const listItems = container.querySelectorAll('li.checklist-item');
        listItems.forEach(li => {
            const item = li.querySelector('input[type="text"]').value;
            const checked = li.querySelector('input[type="checkbox"]').checked;
            items.push({ item, checked });
        });
        return items;
    }

    async function saveTodo(todo) {
        const token = localStorage.getItem('token');
        const method = todo.id ? 'PUT' : 'POST';
        const url = todo.id ? `/todos/${todo.id}` : '/todos';

        const response = await fetch(url, {
            method,
            headers: {
                'Authorization': `${token}`,
                'Content-Type': 'application/json',
            },
            body: JSON.stringify(todo),
        });

        if (response.ok) {
            todosList.innerHTML = '';
            loadTodos();
        }
    }

    async function deleteTodo(id) {
        const token = localStorage.getItem('token');
        await fetch(`/todos/${id}`, {
            method: 'DELETE',
            headers: {
                'Authorization': `${token}`,
            },
        });
        todosList.innerHTML = '';
        loadTodos();
    }

    const logoutBtn = document.getElementById('logoutBtn');
    if (logoutBtn) {
        logoutBtn.addEventListener('click', logout);
    }

    function logout() {
        localStorage.removeItem('token');
        window.location.href = '/app-login';
    }
});
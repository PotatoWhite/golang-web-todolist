(function ($) {
    'use strict';
    $(function () {
        let todoListItem = $('.todo-list');
        let todoListInput = $('.todo-list-input');
        $('.todo-list-add-btn').on("click", function (event) {
            event.preventDefault();
            let item = $(this).prevAll('.todo-list-input').val();

            if (item) {
                $.post("/todos", {name: item}, addItem)
                todoListInput.val("");
            }

        });

        var addItem = function (item) {
            if (item.completed) {
                todoListItem.append("<li class='completed'" + " id='" + item.id + "'><div class='form-check'><label class='form-check-label'><input class='checkbox' type='checkbox' checked='checked' />" + item.name + "<i class='input-helper'></i></label></div><i class='remove mdi mdi-close-circle-outline'></i></li>");
            } else {
                todoListItem.append("<li " + " id='" + item.id + "'><div class='form-check'><label class='form-check-label'><input class='checkbox' type='checkbox' />" + item.name + "<i class='input-helper'></i></label></div><i class='remove mdi mdi-close-circle-outline'></i></li>");
            }
        };

        $.get('/todos', function (items) {
            items.forEach(element => {
                addItem(element)
            })
        })


        todoListItem.on('change', '.checkbox', function () {

            let id = $(this).closest("li").attr('id');
            let $self = $(this);

            let complete = true;
            if ($(this).attr('checked')) {
                complete = false;
            }
            $.get("/complete-todo/" + id + "?complete=" + complete, function (data) {

                if (complete) {
                    $self.attr('checked', 'checked');
                } else {
                    $self.removeAttr('checked');
                }

                $self.closest("li").toggleClass('completed');
            })
        });

        todoListItem.on('click', '.remove', function () {
            // url: todos/id method: DELETE

            let id = $(this).closest("li").attr('id');
            let $self = $(this);
            $.ajax({
                url: "todos/" + id,
                type: "DELETE",
                success: function () {
                    $self.parent().remove();
                }
            })
        });

    });
})(jQuery);
/// <reference path="../pb_data/types.d.ts" />
migrate((db) => {
  const dao = new Dao(db)
  const collection = dao.findCollectionByNameOrId("_pb_users_auth_")

  collection.listRule = "id = @request.auth.id || @request.auth.id ?~ blueprints_via_owner.owner.id"
  collection.viewRule = "id = @request.auth.id || @request.auth.id ?~ blueprints_via_owner.owner.id"

  return dao.saveCollection(collection)
}, (db) => {
  const dao = new Dao(db)
  const collection = dao.findCollectionByNameOrId("_pb_users_auth_")

  collection.listRule = "id = @request.auth.id || @request.auth.id ?~ blueprints_via_owner.users"
  collection.viewRule = "id = @request.auth.id || @request.auth.id ?~ blueprints_via_owner.users"

  return dao.saveCollection(collection)
})

import { models } from '../../wailsjs/go/models'
import { ExportFile } from '../../wailsjs/go/controllers/FileController'

// exportSchiveFile asks the user where to save a copy and writes it there.
// Returns the chosen path, or null when the user cancels the dialog.
export async function exportSchiveFile(file: models.File): Promise<string | null> {
  const destination = await ExportFile(file.id)
  return destination || null
}
